package console

import (
	"fmt"
	"os"
	"pos/internal/domain"
)

func (f *Facade) demo() {
	fmt.Println("=== Демо-сценарий: Преподавательская оценка актёров ===")

	// 1. Создаём студентов
	fmt.Println("\n--- Создаём студентов ---")
	students := []struct{ id, name string }{
		{"jim-carrey", "Джим Керри"},
		{"leonardo-dicaprio", "Леонардо ДиКаприо"},
		{"tom-holland", "Том Холланд"},
		{"tom-cruise", "Том Круз"},
		{"shiloh-buff", "Шайло Бафф"},
		{"eva-elfie", "Ева Елфи"},
		{"keanu-reeves", "Киану Ривз"},
		{"scarlett-johansson", "Скарлетт Йоханссон"},
		{"robert-downey-jr", "Роберт Дауни мл."},
		{"margot-robbie", "Марго Робби"},
		{"brad-pitt", "Брэд Питт"},
		{"johnny-depp", "Джонни Депп"},
		{"zendaya-coleman", "Зендея Коулман"},
		{"timothee-chalamet", "Тимоти Шаламе"},
		{"ryan-gosling", "Райан Гослинг"},
	}

	for _, s := range students {
		if err := f.AddStudent.Execute(s.id, s.name); err != nil {
			fmt.Printf("  Ошибка добавления %s: %v\n", s.name, err)
		} else {
			fmt.Printf("  ✓ %s\n", s.name)
		}
	}

	// 2. Создаём критерии
	fmt.Println("\n--- Создаём критерии ---")
	criteria := []struct{ name string; maxScore int; weight float64 }{
		{"Харизма на экране", 10, 0.3},
		{"Умение играть драму", 10, 0.2},
		{"Комедийный талант", 10, 0.2},
		{"Физическая подготовка", 10, 0.15},
		{"Работа с импровизацией", 10, 0.15},
	}

	for _, c := range criteria {
		id, err := f.AddCriterion.Execute(c.name, c.maxScore, c.weight)
		if err != nil {
			fmt.Printf("  Ошибка добавления критерия '%s': %v\n", c.name, err)
		} else {
			fmt.Printf("  ✓ %s (id=%s, max=%d, weight=%.2f)\n", c.name, id, c.maxScore, c.weight)
		}
	}

	// 3. Выставляем оценки (с пересдачами для демонстрации стратегий)
	fmt.Println("\n--- Выставляем оценки ---")

	// Джим Керри: отличная комедия с пересдачей (улучшил), но слабая драма
	f.AddEvaluation.Execute("jim-carrey", "комедийный-талант", 9, "Маска легла идеально")
	f.AddEvaluation.Execute("jim-carrey", "комедийный-талант", 10, "Пересдача: теперь и лицо, и голос")
	f.AddEvaluation.Execute("jim-carrey", "харизма-на-экране", 10, "Заполняет кадр собой")
	f.AddEvaluation.Execute("jim-carrey", "умение-играть-драму", 4, "Серьёзные роли пока не его")
	f.AddEvaluation.Execute("jim-carrey", "работа-с-импровизацией", 10, "Импровизирует лучше, чем по сценарию")

	// Леонардо ДиКаприо: отличная драма, но нет оценок по комедии и импровизации
	f.AddEvaluation.Execute("leonardo-dicaprio", "умение-играть-драму", 10, "Медведь в лесу — это был подвиг")
	f.AddEvaluation.Execute("leonardo-dicaprio", "харизма-на-экране", 9, "Взгляд, за который дают Оскар")
	f.AddEvaluation.Execute("leonardo-dicaprio", "физическая-подготовка", 8, "Выжил в холоде и грязи")
	f.AddEvaluation.Execute("leonardo-dicaprio", "комедийный-талант", 5, "Смешно, но не специально")

	// Том Холланд: пересдача драмы (улучшил с 6 до 7)
	f.AddEvaluation.Execute("tom-holland", "физическая-подготовка", 10, "Кувырки и паркур без дублёра")
	f.AddEvaluation.Execute("tom-holland", "комедийный-талант", 8, "Обаятельный неловкий герой")
	f.AddEvaluation.Execute("tom-holland", "умение-играть-драму", 6, "Первый блин комом")
	f.AddEvaluation.Execute("tom-holland", "умение-играть-драму", 7, "Пересдача: стал глубже")
	f.AddEvaluation.Execute("tom-holland", "харизма-на-экране", 8, "Соседский дружелюбный тон")

	// Том Круз: все критерии, но без пересдач
	f.AddEvaluation.Execute("tom-cruise", "физическая-подготовка", 10, "Висел на самолёте сам")
	f.AddEvaluation.Execute("tom-cruise", "харизма-на-экране", 9, "Бежит так, что зал встаёт")
	f.AddEvaluation.Execute("tom-cruise", "умение-играть-драму", 7, "Тихие сцены даются сложнее экшена")
	f.AddEvaluation.Execute("tom-cruise", "комедийный-талант", 6, "Серьёзный в комедии")
	f.AddEvaluation.Execute("tom-cruise", "работа-с-импровизацией", 5, "Любит точный сценарий")

	// Шайло Бафф: пересдача драмы (улучшил с 5 до 7)
	f.AddEvaluation.Execute("shiloh-buff", "харизма-на-экране", 6, "Есть искра, но пока тлеет")
	f.AddEvaluation.Execute("shiloh-buff", "умение-играть-драму", 5, "Слёзы были, чувств не было")
	f.AddEvaluation.Execute("shiloh-buff", "умение-играть-драму", 7, "Пересдача: теперь плачет по делу")
	f.AddEvaluation.Execute("shiloh-buff", "комедийный-талант", 8, "Тайминг шутки на уровне")
	f.AddEvaluation.Execute("shiloh-buff", "физическая-подготовка", 4, "Прыжок с дивана — предел")
	f.AddEvaluation.Execute("shiloh-buff", "работа-с-импровизацией", 6, "Импровизирует робко")

	// Ева Елфи: нет оценок по физподготовке
	f.AddEvaluation.Execute("eva-elfie", "комедийный-талант", 9, "Сарказм как оружие")
	f.AddEvaluation.Execute("eva-elfie", "работа-с-импровизацией", 8, "Держит сцену без текста")
	f.AddEvaluation.Execute("eva-elfie", "харизма-на-экране", 7, "Харизма с характером")
	f.AddEvaluation.Execute("eva-elfie", "умение-играть-драму", 6, "Драма пока на втором плане")

	// Киану Ривз: все критерии
	f.AddEvaluation.Execute("keanu-reeves", "физическая-подготовка", 9, "Кунг-фу и огнестрел без каскадёра")
	f.AddEvaluation.Execute("keanu-reeves", "умение-играть-драму", 8, "Грустный взгляд работает")
	f.AddEvaluation.Execute("keanu-reeves", "харизма-на-экране", 9, "Спокойная сила")
	f.AddEvaluation.Execute("keanu-reeves", "комедийный-талант", 5, "Смех редкий, но искренний")
	f.AddEvaluation.Execute("keanu-reeves", "работа-с-импровизацией", 7, "Импровизирует с достоинством")

	// Скарлетт Йоханссон: все критерии
	f.AddEvaluation.Execute("scarlett-johansson", "харизма-на-экране", 10, "Голос и взгляд — магнит")
	f.AddEvaluation.Execute("scarlett-johansson", "умение-играть-драму", 8, "Тихая драма без надрыва")
	f.AddEvaluation.Execute("scarlett-johansson", "комедийный-талант", 7, "Ирония в тему")
	f.AddEvaluation.Execute("scarlett-johansson", "физическая-подготовка", 8, "Экипировка не помеха")
	f.AddEvaluation.Execute("scarlett-johansson", "работа-с-импровизацией", 6, "Предпочитает текст")

	// Роберт Дауни мл.: пересдача драмы (улучшил с 7 до 8)
	f.AddEvaluation.Execute("robert-downey-jr", "работа-с-импровизацией", 10, "Сценаристы плачут от его импровизаций")
	f.AddEvaluation.Execute("robert-downey-jr", "комедийный-талант", 9, "Шутки рождаются на ходу")
	f.AddEvaluation.Execute("robert-downey-jr", "харизма-на-экране", 9, "Обаяние через край")
	f.AddEvaluation.Execute("robert-downey-jr", "умение-играть-драму", 7, "Драма сквозь улыбку")
	f.AddEvaluation.Execute("robert-downey-jr", "умение-играть-драму", 8, "Пересдача: добавил тишины")
	f.AddEvaluation.Execute("robert-downey-jr", "физическая-подготовка", 6, "Костюм тяжелее, чем кажется")

	// Марго Робби: все критерии
	f.AddEvaluation.Execute("margot-robbie", "харизма-на-экране", 9, "Зажигает с первой сцены")
	f.AddEvaluation.Execute("margot-robbie", "комедийный-талант", 8, "Дерзкая комедия")
	f.AddEvaluation.Execute("margot-robbie", "умение-играть-драму", 7, "Сдержанно, но верно")
	f.AddEvaluation.Execute("margot-robbie", "физическая-подготовка", 8, "Коньки и прыжки — норма")
	f.AddEvaluation.Execute("margot-robbie", "работа-с-импровизацией", 7, "Ловко выходит из пауз")

	// Брэд Питт: все критерии
	f.AddEvaluation.Execute("brad-pitt", "харизма-на-экране", 9, "Расслабленная уверенность")
	f.AddEvaluation.Execute("brad-pitt", "умение-играть-драму", 8, "Усталость и глубина")
	f.AddEvaluation.Execute("brad-pitt", "комедийный-талант", 6, "Комедия не его жанр")
	f.AddEvaluation.Execute("brad-pitt", "физическая-подготовка", 7, "Драки держит")
	f.AddEvaluation.Execute("brad-pitt", "работа-с-импровизацией", 5, "Играет по битам")

	// Джонни Депп: все критерии
	f.AddEvaluation.Execute("johnny-depp", "работа-с-импровизацией", 10, "Приносит реквизит из ниоткуда")
	f.AddEvaluation.Execute("johnny-depp", "комедийный-талант", 9, "Эксцентрика — его стихия")
	f.AddEvaluation.Execute("johnny-depp", "умение-играть-драму", 8, "Меланхолия в кадре")
	f.AddEvaluation.Execute("johnny-depp", "харизма-на-экране", 9, "Странный и притягательный")
	f.AddEvaluation.Execute("johnny-depp", "физическая-подготовка", 5, "Походка важнее мышц")

	// Зендея Коулман: все критерии
	f.AddEvaluation.Execute("zendaya-coleman", "умение-играть-драму", 9, "Эмоция без слов")
	f.AddEvaluation.Execute("zendaya-coleman", "харизма-на-экране", 9, "Современная икона")
	f.AddEvaluation.Execute("zendaya-coleman", "комедийный-талант", 7, "Лёгкая ирония")
	f.AddEvaluation.Execute("zendaya-coleman", "физическая-подготовка", 8, "Танец как база")
	f.AddEvaluation.Execute("zendaya-coleman", "работа-с-импровизацией", 7, "Импровизирует уверенно")

	// Тимоти Шаламе: все критерии
	f.AddEvaluation.Execute("timothee-chalamet", "умение-играть-драму", 9, "Хрупкость и надрыв")
	f.AddEvaluation.Execute("timothee-chalamet", "комедийный-талант", 7, "Самоирония есть")
	f.AddEvaluation.Execute("timothee-chalamet", "харизма-на-экране", 8, "Скандальная харизма")
	f.AddEvaluation.Execute("timothee-chalamet", "физическая-подготовка", 6, "Танец лучше бега")
	f.AddEvaluation.Execute("timothee-chalamet", "работа-с-импровизацией", 7, "Импровизирует в ритме")

	// Райан Гослинг: все критерии
	f.AddEvaluation.Execute("ryan-gosling", "харизма-на-экране", 10, "Стиль важнее реплик")
	f.AddEvaluation.Execute("ryan-gosling", "умение-играть-драму", 8, "Тихая боль в кадре")
	f.AddEvaluation.Execute("ryan-gosling", "комедийный-талант", 7, "Сухой юмор")
	f.AddEvaluation.Execute("ryan-gosling", "физическая-подготовка", 7, "Пианино и бокс")
	f.AddEvaluation.Execute("ryan-gosling", "работа-с-импровизацией", 6, "Держит паузу")

	fmt.Println("  ✓ Оценки выставлены")

	// 4. Выводим списки
	fmt.Println("\n--- Студенты ---")
	studentsList, _ := f.ListStudents.Execute()
	for _, s := range studentsList {
		fmt.Printf("  %s — %s (оценок: %d)\n", s.ID().String(), s.Name(), len(s.Evaluations()))
	}

	fmt.Println("\n--- Критерии ---")
	f.printCriteria()

	// 5. Рассчитываем итоги для нескольких студентов по обеим стратегиям
	fmt.Println("\n--- Расчёт итогов (демонстрация стратегий) ---")
	last := domain.LastAttemptPolicy{}
	highest := domain.HighestScorePolicy{}

	testStudents := []string{"jim-carrey", "tom-holland", "shiloh-buff", "robert-downey-jr"}

	for _, sid := range testStudents {
		fmt.Printf("\n%s:\n", sid)

		gLast, errL := f.CalculateFinalGrade.Execute(sid, last)
		gHigh, errH := f.CalculateFinalGrade.Execute(sid, highest)

		if errL == nil {
			fmt.Printf("  last (последняя попытка):    %s\n", gLast.String())
		} else {
			fmt.Printf("  last: ошибка: %v\n", errL)
		}

		if errH == nil {
			fmt.Printf("  highest (лучшая попытка):    %s\n", gHigh.String())
		} else {
			fmt.Printf("  highest: ошибка: %v\n", errH)
		}

		if errL == nil && errH == nil && gLast.Value() != gHigh.Value() {
			fmt.Printf("  ⚡ Разница: %.2f (стратегии дают разные результаты!)\n",
				gHigh.Value()-gLast.Value())
		}
	}

	// 6. Экспортируем итоговую ведомость
	fmt.Println("\n--- Экспорт итоговой ведомости ---")
	data, err := f.ExportFinalGrades.Execute(last)
	if err == nil {
		path := "./data/final_grades_demo.csv"
		if err := os.WriteFile(path, data, 0o644); err != nil {
			fmt.Println("Ошибка записи файла:", err)
		} else {
			fmt.Printf("✓ Экспортировано в %s (стратегия: last)\n", path)
		}
	} else {
		fmt.Println("Ошибка экспорта:", err)
	}

	fmt.Println("\n=== Демо-сценарий завершён ===")
}