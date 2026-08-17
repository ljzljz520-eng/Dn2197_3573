package fixtures

import "gradebook/domain"

func Record(id, name, program string, scores ...float64) domain.Gradebook {
	courses := []domain.CourseGrade{
		domain.NewCourseGrade("GRA501", "Research Methods", 3, scoreAt(scores, 0), ""),
		domain.NewCourseGrade("GRA502", "Advanced Theory", 3, scoreAt(scores, 1), ""),
		domain.NewCourseGrade("GRA503", "Scholarly Practice", 2, scoreAt(scores, 2), ""),
	}
	return domain.Gradebook{RecordID: id, Student: domain.NewStudent(id, name, program, 2025, id+"@example.edu", "Dr. Mentor"), DegreeCourses: courses}
}

func scoreAt(scores []float64, index int) float64 {
	if index < len(scores) {
		return scores[index]
	}
	return 70
}

func Seminar(code, name, letter string, score float64) domain.CourseGrade {
	return domain.NewCourseGrade(code, name, 2, score, letter)
}

func Batch() []domain.Gradebook {
	return []domain.Gradebook{Record("S001", "Ada Lovelace", "Computing", 96, 91, 89), Record("S002", "Grace Hopper", "Computing", 88, 84, 86), Record("S003", "Katherine Johnson", "Statistics", 78, 82, 80)}
}

func WithSeminar(record domain.Gradebook, seminar domain.CourseGrade) domain.Gradebook {
	record.Seminars = append(record.Seminars, seminar)
	return record
}
