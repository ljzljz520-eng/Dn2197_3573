package grading

import (
	"fmt"
	"gradebook/domain"
	"sort"
	"strings"
)

type TranscriptLine struct {
	CourseCode string
	CourseName string
	Credits    int
	Score      float64
	Letter     string
	Band       string
}
type Transcript struct {
	StudentID   string
	StudentName string
	Lines       []TranscriptLine
	Average     float64
	Standing    string
}

func BuildTranscript(record domain.Gradebook) Transcript {
	transcript := Transcript{StudentID: record.Student.ID, StudentName: record.Student.Name, Lines: []TranscriptLine{}}
	for _, course := range record.DegreeCourses {
		transcript.Lines = append(transcript.Lines, lineFor(course, false))
	}
	for _, course := range record.Seminars {
		transcript.Lines = append(transcript.Lines, lineFor(course, true))
	}
	transcript.Average = WeightedAverage(record)
	transcript.Standing = Standing(transcript.Average, record.IsComplete())
	sort.Slice(transcript.Lines, func(i, j int) bool { return transcript.Lines[i].CourseCode < transcript.Lines[j].CourseCode })
	return transcript
}

func lineFor(course domain.CourseGrade, seminar bool) TranscriptLine {
	score := course.Score
	if seminar && course.Letter != "" {
		score = LetterScore(course.Letter)
	}
	return TranscriptLine{CourseCode: course.CourseCode, CourseName: course.CourseName, Credits: course.Credits, Score: score, Letter: course.Letter, Band: PercentageBand(score)}
}

func Standing(average float64, complete bool) string {
	if !complete {
		return "incomplete"
	}
	switch {
	case average >= 90:
		return "honors"
	case average >= 80:
		return "good standing"
	case average >= 60:
		return "satisfactory"
	default:
		return "academic review"
	}
}

func RenderTranscript(transcript Transcript) string {
	lines := []string{fmt.Sprintf("%s (%s)", transcript.StudentName, transcript.StudentID)}
	for _, line := range transcript.Lines {
		lines = append(lines, fmt.Sprintf("%s %s %d %.2f %s", line.CourseCode, line.CourseName, line.Credits, line.Score, line.Band))
	}
	lines = append(lines, fmt.Sprintf("average %.2f; %s", transcript.Average, transcript.Standing))
	return strings.Join(lines, "\n")
}

func CreditsCompleted(record domain.Gradebook) int {
	total := 0
	for _, course := range record.AllCourses() {
		score := course.Score
		if course.Letter != "" {
			score = LetterScore(course.Letter)
		}
		if score >= 60 {
			total += course.Credits
		}
	}
	return total
}

func EligibleForGraduation(record domain.Gradebook, minimumCredits int) bool {
	return record.IsComplete() && CreditsCompleted(record) >= minimumCredits && WeightedAverage(record) >= 60
}

func CourseOutcome(course domain.CourseGrade) string {
	score := course.Score
	if course.Letter != "" {
		score = LetterScore(course.Letter)
	}
	if score >= 60 {
		return "credit"
	}
	return "no credit"
}
