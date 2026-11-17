package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"pos/internal/domain"
)

var stdin = bufio.NewReader(os.Stdin)

func readLine() string {
	line, _ := stdin.ReadString('\n')
	return strings.TrimSpace(line)
}

func readInt(prompt string) (int, error) {
	fmt.Print(prompt)
	return strconv.Atoi(readLine())
}

func readFloat(prompt string) (float64, error) {
	fmt.Print(prompt)
	return strconv.ParseFloat(readLine(), 64)
}

func (f *Facade) choosePolicy() domain.AttemptPolicy {
	fmt.Println("  Выберите политику пересдач:")
	fmt.Println("    1. последняя попытка")
	fmt.Println("    2. лучшая попытка")
	fmt.Print("  > ")
	switch readLine() {
	case "2":
		return domain.HighestScorePolicy{}
	default:
		return domain.LastAttemptPolicy{}
	}
}

func (f *Facade) Run() {
	fmt.Println("=== Модуль ПОС: Преподавательская Оценка Студентов ===")
	for {
		fmt.Println()
		fmt.Println("Меню:")
		fmt.Println("  1. Добавить студента")
		fmt.Println("  2. Список студентов")
		fmt.Println("  3. Добавить критерий")
		fmt.Println("  4. Список критериев")
		fmt.Println("  5. Удалить критерий")
		fmt.Println("  6. Выставить оценку")
		fmt.Println("  7. Оценки студента")
		fmt.Println("  8. Исправить оценку")
		fmt.Println("  9. Удалить оценку")
		fmt.Println(" 10. Рассчитать итог по студенту")
		fmt.Println(" 11. Экспорт итогов в CSV")
		fmt.Println(" 12. Демо-сценарий")
		fmt.Println("  0. Выход")
		fmt.Print("> ")

		switch readLine() {
		case "1":
			fmt.Print("ID студента: ")
			id := readLine()
			fmt.Print("Имя: ")
			name := readLine()
			if err := f.AddStudent.Execute(id, name); err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Студент добавлен")
			}
		case "2":
			students, err := f.ListStudents.Execute()
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			if len(students) == 0 {
				fmt.Println("Студентов нет")
				continue
			}
			for _, s := range students {
				fmt.Printf("  %s — %s (оценок: %d)\n", s.ID().String(), s.Name(), len(s.Evaluations()))
			}
		case "3":
			fmt.Print("Название критерия: ")
			name := readLine()
			maxScore, err := readInt("Максимум: ")
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			weight, err := readFloat("Вес (например 0.3): ")
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			id, err := f.AddCriterion.Execute(name, maxScore, weight)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Критерий добавлен, id =", id)
			}
		case "4":
			f.printCriteria()
		case "5":
			fmt.Print("ID критерия: ")
			if err := f.RemoveCriterion.Execute(readLine()); err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Удалено")
			}
		case "6":
			fmt.Print("ID студента: ")
			sid := readLine()
			fmt.Print("ID критерия: ")
			cid := readLine()
			score, err := readInt("Балл: ")
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			fmt.Print("Комментарий (опционально): ")
			comment := readLine()
			eid, err := f.AddEvaluation.Execute(sid, cid, score, comment)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Оценка выставлена, id =", eid)
			}
		case "7":
			fmt.Print("ID студента: ")
			f.printEvaluations(readLine())
		case "8":
			fmt.Print("ID студента: ")
			sid := readLine()
			fmt.Print("ID оценки: ")
			eid := readLine()
			score, err := readInt("Новый балл: ")
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			if err := f.CorrectEvaluation.Execute(sid, eid, score); err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Исправлено")
			}
		case "9":
			fmt.Print("ID студента: ")
			sid := readLine()
			fmt.Print("ID оценки: ")
			eid := readLine()
			if err := f.RemoveEvaluation.Execute(sid, eid); err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Удалено")
			}
		case "10":
			fmt.Print("ID студента: ")
			sid := readLine()
			policy := f.choosePolicy()
			grade, err := f.CalculateFinalGrade.Execute(sid, policy)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Printf("Итог (%s): %s\n", policy.Name(), grade.String())
			}
		case "11":
			policy := f.choosePolicy()
			data, err := f.ExportFinalGrades.Execute(policy)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			path := "./data/final_grades.csv"
			if err := os.WriteFile(path, data, 0o644); err != nil {
				fmt.Println("Ошибка записи:", err)
				continue
			}
			fmt.Println("Экспортировано в", path)
		case "12":
			f.demo()
		case "0":
			return
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}
