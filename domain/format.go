package domain

import (
	"fmt"
	"sort"
	"strings"
)

func Normalize(g *Gradebook) {
	g.Student.ID = strings.TrimSpace(g.Student.ID)
	g.Student.Name = strings.TrimSpace(g.Student.Name)
	g.Student.Program = strings.TrimSpace(g.Student.Program)
	g.Student.Email = strings.TrimSpace(strings.ToLower(g.Student.Email))
	g.Student.Advisor = strings.TrimSpace(g.Student.Advisor)
	for i := range g.DegreeCourses {
		g.DegreeCourses[i].CourseCode = strings.ToUpper(strings.TrimSpace(g.DegreeCourses[i].CourseCode))
		g.DegreeCourses[i].CourseName = strings.TrimSpace(g.DegreeCourses[i].CourseName)
	}
	for i := range g.Seminars {
		g.Seminars[i].CourseCode = strings.ToUpper(strings.TrimSpace(g.Seminars[i].CourseCode))
		g.Seminars[i].CourseName = strings.TrimSpace(g.Seminars[i].CourseName)
		g.Seminars[i].Letter = strings.ToUpper(strings.TrimSpace(g.Seminars[i].Letter))
	}
}

func CourseCodes(g Gradebook) []string {
	codes := make([]string, 0, len(g.AllCourses()))
	for _, c := range g.AllCourses() {
		codes = append(codes, c.CourseCode)
	}
	sort.Strings(codes)
	return codes
}

func Summary(g Gradebook) string {
	return fmt.Sprintf("%s | %s | average %.2f | confirmations %d", g.DisplayName(), strings.Join(CourseCodes(g), ","), g.Average, g.ConfirmationCount)
}
