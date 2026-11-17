# Модуль «ПОС» (Преподавательская Оценка Студентов)

ФИО: Наконечный Артемий Алексеевич
---

## Архитектура

Выбрана архитектура **Ports and Adapters (Hexagonal)**. Она обеспечивает чёткое разделение ответственности и независимость ядра бизнес-логики от фреймворков, UI и механизмов хранения.

### Структура пакетов

```text
pos/
├── cmd/pos/main.go                # Composition Root (сборка DI)
├── go.mod
└── internal/
    ├── domain/                    # Ядро: Entities, VO, Domain Services, ошибки
    │   ├── vo/                    # Value Objects (Score, Weight, IDs и т.д.)
    │   ├── student.go             # Агрегат Student
    │   ├── rubric.go              # Агрегат Rubric
    │   ├── criterion.go           # Entity внутри Rubric
    │   ├── evaluation.go          # Entity внутри Student
    │   ├── attempt_policy.go      # Strategy-интерфейс + реализации
    │   ├── final_grade_calculator.go  # Domain Service
    │   └── errors.go              # Доменные ошибки
    ├── application/               # Сценарии использования и порты
    │   ├── port/                  # Интерфейсы репозиториев и экспортеров
    │   ├── dto/                   # DTO для передачи данных между слоями
    │   └── usecase/               # Use Cases (AddStudent, CalculateFinalGrade и др.)
    └── infrastructure/            # Адаптеры
        ├── persistence/           # JSON-репозитории (Student, Rubric)
        ├── export/                # CSV-экспортер
        └── console/               # CLI-меню и Facade
```

### Направление зависимостей

Зависимости направлены строго **внутрь**, к домену:

```text
                    ┌──────────────────────┐
                    │      cmd/pos         │ (Composition Root)
                    └──────────┬───────────┘
                               │ зависит от всего
              ┌────────────────┼────────────────┐
              ▼                ▼                ▼
      ┌───────────────┐ ┌──────────────┐ ┌──────────────────┐
      │ infrastructure│ │ application  │ │    domain        │
      │  (adapters)   │ │ (use cases,  │ │ (entities, VO,   │
      │               │ │   ports)     │ │  services)       │
      └──────┬────────┘ └──────┬───────┘ └────────┬─────────┘
             │                 │                  ▲
             │ реализует       │ зависит от       │
             │ порты           │ domain           │
             └─────────────────┼──────────────────┘
                               ▼
```

`domain` не импортирует ничего, кроме стандартной библиотеки Go. `infrastructure` зависит от `application/port`, а не наоборот. `cmd/pos/main.go` — единственное место, где собирается граф зависимостей.

---

## Domain-Driven Design 

### 1. Ubiquitous Language

| В ТЗ | В коде |
|---|---|
| студент | `domain.Student`, `vo.StudentID` |
| критерий | `domain.Criterion`, `vo.CriterionID` |
| оценка | `domain.Evaluation`, `vo.EvaluationID` |
| итог | `vo.FinalGrade`, `domain.FinalGradeCalculator` |
| рубрикатор | `domain.Rubric` |
| политика пересдач | `domain.AttemptPolicy` |


### 2. Entity и Value Object
Entity:
- `Student` — идентифицируется по `StudentID`, может накапливать `Evaluation`.
- `Criterion` — идентифицируется по `CriterionID`, имеет имя, максимум и вес.
- `Evaluation` — идентифицируется по `EvaluationID`, имеет номер попытки, балл, дату, комментарий.

VO:
- `Score` — целое неотрицательное число.
- `MaxScore` — строго положительное целое.
- `Weight` — строго положительное вещественное число.
- `FinalGrade` — вещественное число в диапазоне `[0, 100]`.
- `StudentID`, `CriterionID`, `EvaluationID` — непустые строки-обёртки.

Все VO создаются только через валидирующие конструкторы (`NewScore`, `NewMaxScore`, `NewWeight`, `NewFinalGrade`, `NewStudentID` и т.д.), что гарантирует отсутствие невалидных значений в системе.

### 3. Инварианты

| Инвариант | Где охраняется |
|---|---|
| `score ≤ maxScore` | `newEvaluation` через `score.IsWithin(maxScore)` |
| `score ≥ 0` | `vo.NewScore` |
| `maxScore > 0` | `vo.NewMaxScore` |
| `weight > 0` | `vo.NewWeight` |
| `finalGrade ∈ [0, 100]` | `vo.NewFinalGrade` |
| Оценка только по существующему критерию | `usecase.AddEvaluation` через `rubric.FindCriterion` |
| Нельзя удалить критерий с оценками | `usecase.RemoveCriterion` перебором всех студентов |
| Нельзя удалить несуществующую оценку | `Student.RemoveEvaluation` |
| Имя студента не пустое | `domain.NewStudent` |
| Имя критерия не пустое | `domain.NewCriterion` |
| ID не пустые | конструкторы VO |
| Итог согласован с оценками | `FinalGradeCalculator.Calculate` |

### 4. Repository и Domain Service

**Repository Ports** (`StudentRepository`, `RubricRepository`) объявлены в `application/port`. Домен не знает, как сохраняются данные (JSON, БД, in-memory), он только использует интерфейсы.

Логика расчёта итогового балла не принадлежит ни `Student`, ни `Rubric` по отдельности. Это операция, координирующая оба агрегата и применяющая `AttemptPolicy`. Формула:

```
final = Σ(score / maxScore × weight) / Σ(weight) × 100
```

---

## Паттерны GoF

### Порождающий: Factory 

Конструкторы `NewStudent`, `NewCriterion`, `newEvaluation`, `NewScore`, `NewMaxScore`, `NewWeight`, `NewFinalGrade`, `NewStudentID` и т.д.

Гарантирует, что в систему не попадёт невалидный объект. Инварианты проверяются в момент создания, предотвращая появление анемичной модели. Например, невозможно создать `Score` с отрицательным значением или `Evaluation`, где `score > maxScore`.

### Структурный: Facade 

`console.Facade` в пакете `infrastructure/console`.

Скрывает от CLI-меню сложность оркестрации множества use-case'ов. Меню вызывает один метод (например, `f.AddStudent.Execute()`), не зная о репозиториях, маппинге DTO и деталях домена. Вся зависимость UI от приложения проходит через единую точку входа.

### Поведенческий: Strategy 

Интерфейс `domain.AttemptPolicy` с реализациями `LastAttemptPolicy` и `HighestScorePolicy`.

Позволяет менять правило пересдач (брать последнюю или лучшую попытку) без изменения кода `FinalGradeCalculator` и без создания `if/else` спагетти. Пользователь выбирает политику в меню, и она передаётся в Domain Service как параметр.

---

##  Dependency Injection (DI)


```go
func main() {
    dataDir := "./data"
    studentRepo := persistence.NewJSONStudentRepo(filepath.Join(dataDir, "students.json"))
    rubricRepo  := persistence.NewJSONRubricRepo(filepath.Join(dataDir, "rubric.json"))
    exporter    := export.NewCSVExporter()
    calc        := domain.NewFinalGradeCalculator()

    deps := console.Deps{
        AddStudent:          usecase.NewAddStudent(studentRepo),
        ListStudents:        usecase.NewListStudents(studentRepo),
        AddCriterion:        usecase.NewAddCriterion(rubricRepo),
        // ...
        CalculateFinalGrade: usecase.NewCalculateFinalGrade(studentRepo, rubricRepo, calc),
        ExportFinalGrades:   usecase.NewExportFinalGrades(studentRepo, rubricRepo, calc, exporter),
    }
    console.NewFacade(deps).Run()
}
```

Не требует тяжёлых фреймворков, легко читается, быстро работает и позволяет легко подменять реализации (например, на моки) в юнит-тестах.

---

## Тестирование

Реализованы юнит-тесты на инварианты и Domain Service с использованием стандартных пакетов

**Покрытие:**

- **Value Objects:** валидация границ (`Score < 0`, `Weight <= 0`, `FinalGrade` вне `[0, 100]`, пустые ID).
- **Score.IsWithin:** граница `≤ max`, выход за `max`.
- **Entities:** автоинкремент попыток (`attempt`), запрет на превышение `maxScore` при добавлении и исправлении оценки, удаление несуществующей оценки.
- **Rubric:** slug-нормализация имени, дубликаты, пустое имя, удаление.
- **AttemptPolicy:** выбор `last` / `highest`, пустой список.
- **FinalGradeCalculator:** нет критериев, нет оценок, `last` vs `highest`, взвешенное среднее.

**Запуск тестов:**

```bash
go test ./internal/domain/... -v
```

---

## Запуск 

```bash
go run cmd/pos/main.go
```

Опционально можно указать папку для данных:

```bash
go run cmd/pos/main.go ./my_data
```

### Демонстрационный сценарий

В главном меню выберите пункт **`12. Демо-сценарий`**. Приложение автоматически:

1. Создаст двух студентов (`demo-s1`, `demo-s2`) и три критерия (`Homework` 0.3, `Exam` 0.5, `Attendance` 0.2).
2. Выставит оценки, включая пересдачу с ухудшением результата для `demo-s2`.
3. Выведет списки студентов, критериев и оценок.
4. Рассчитает итоговый балл по обеим стратегиям (`last` и `highest`).
5. Экспортирует итоговую ведомость в файл `./data/final_grades_demo.csv`.

---

## Роль ИИ в решении

Искусственный интеллект:
1. **DDD:** Консультировал в моментах по типу строгости домена, рамок контекста, что можно брать за отдельные VO `StudentID` `CriterionID`, а где можно обойтись и без них `StudentName` `CriterionName` `Comment`
2. **Помощь в работе с json** 
3. **Написание тестов для** `FinalGradeCalculator` и `AttemptPolicy`, покрывающую граничные условия, которые легко упустить при ручном написании.
4. **Генерация шаблона для README**