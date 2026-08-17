package fixtures

import "gradebook/domain"

type Scenario struct {
	Name     string
	Record   domain.Gradebook
	Operator string
}

func Scenarios() []Scenario {
	return []Scenario{
		{Name: "honors", Record: WithSeminar(Record("S101", "Lin", "Physics", 98, 95, 93), Seminar("SEM601", "Research Seminar", "A+", 0)), Operator: "fixture"},
		{Name: "merit", Record: WithSeminar(Record("S102", "Noor", "Physics", 86, 84, 82), Seminar("SEM601", "Research Seminar", "B+", 0)), Operator: "fixture"},
		{Name: "review", Record: Record("S103", "Mina", "History", 55, 61, 58), Operator: "fixture"},
	}
}

func ScenarioMap() map[string]Scenario {
	result := map[string]Scenario{}
	for _, scenario := range Scenarios() {
		result[scenario.Name] = scenario
	}
	return result
}

func SeminarRecord(id string) domain.Gradebook {
	return WithSeminar(Record(id, "Seminar Student", "Interdisciplinary", 75, 80, 85), Seminar("SEM601", "Research Seminar", "B", 0))
}

func InvalidRecord() domain.Gradebook {
	return domain.Gradebook{RecordID: "bad", Student: domain.NewStudent("", "", "", 1900, "", "")}
}
