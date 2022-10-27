package data

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Element struct {
	id      int
	title   string
	details string
}

type Question struct {
	id       int
	text     string
	isBinary bool
}

type Answer struct {
	id       int
	element  int
	question int
	answer   int
}

type Data struct {
	db        *sql.DB
	elements  []Element
	questions []Question
}

func getElements(db *sql.DB) ([]Element, error) {
	rows, err := db.Query("SELECT id, title, details FROM element WHERE active = TRUE")
	if err != nil {
		return nil, fmt.Errorf("could not fetch elements from DB: %w", err)
	}

	elems := make([]Element, 0)

	for rows.Next() {
		e := Element{}
		if err = rows.Scan(&e.id, &e.title, &e.details); err != nil {
			return nil, fmt.Errorf("could not read element from DB: %w", err)
		}

		elems = append(elems, e)
	}

	if err = rows.Close(); err != nil {
		return nil, fmt.Errorf("could not close result of query: %w", err)
	}

	return elems, nil
}

func getQuestions(db *sql.DB) ([]Question, error) {
	rows, err := db.Query("SELECT id, text, isBinary FROM question WHERE active = TRUE")
	if err != nil {
		return nil, fmt.Errorf("could not fetch questions from DB: %w", err)
	}

	questions := make([]Question, 0)

	for rows.Next() {
		q := Question{}
		if err = rows.Scan(&q.id, &q.text, &q.isBinary); err != nil {
			return nil, fmt.Errorf("could not read question from DB: %w", err)
		}

		questions = append(questions, q)
	}

	if err = rows.Close(); err != nil {
		return nil, fmt.Errorf("could not close result of query: %w", err)
	}

	return questions, nil
}

func New() Data {
	// Initialize random seed, though not sure if this qualifies as too "magic"
	rand.Seed(time.Now().Unix())

	// Connect to the DB
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(fmt.Errorf("could not open connection to DB: %w", err))
	}

	elements, err := getElements(db)
	if err != nil {
		panic(fmt.Errorf("could not get elements: %w", err))
	}

	questions, err := getQuestions(db)
	if err != nil {
		panic(fmt.Errorf("could not get questions: %w", err))
	}

	return Data{
		db,
		elements,
		questions,
	}
}

func CreateDB() error {
	// Connect to the DB
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return fmt.Errorf("could not open connection to DB: %w", err)
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE `question` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`text` text NOT NULL," +
		"`isBinary` int NOT NULL," +
		"`active` int DEFAULT TRUE NOT NULL," +
		")")
	if err != nil {
		return fmt.Errorf("could not create question table: %w", err)
	}

	_, err = db.Exec("CREATE TABLE `element` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`title` text NOT NULL," +
		"`details` text NULL" +
		"`active` int DEFAULT TRUE NOT NULL," +
		")")
	if err != nil {
		return fmt.Errorf("could not create element table: %w", err)
	}

	_, err = db.Exec("CREATE TABLE `answer` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`element` int NOT NULL," +
		"`question` int NOT NULL," +
		"`answer` int NOT NULL, " +
		"`created` TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
		"FOREIGN KEY(`element`) REFERENCES `element`(`id`)," +
		"FOREIGN KEY(`question`) REFERENCES `question`(`id`)" +
		")")
	if err != nil {
		return fmt.Errorf("could not create answer table: %w", err)
	}

	return nil
}

type Prompt struct {
	question  Question
	element   Element
	startTime time.Time
}

type Renderable struct {
	ElementTitle     string
	ElementDetails   string
	QuestionText     string
	QuestionIsBinary bool
}

func (p Prompt) Get() Renderable {
	return Renderable{
		ElementTitle:     p.element.title,
		ElementDetails:   p.element.details,
		QuestionText:     p.question.text,
		QuestionIsBinary: p.question.isBinary,
	}
}

func (d Data) Ask() Prompt {
	question := d.questions[rand.Intn(len(d.questions))]
	element := d.elements[rand.Intn(len(d.elements))]

	return Prompt{
		question,
		element,
		time.Now(),
	}
}

// deleteLastIfMatch checks if the last saved answer matches the current prompt, deleting it if so
func deleteLastIfMatch(tx *sql.Tx, prompt Prompt) error {
	stmt, err := tx.Prepare("delete from answer where element = ? and question = ? and id = (" +
		"select id from answer order by created desc limit 1);")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	if _, err := stmt.Exec(prompt.element.id, prompt.question.id); err != nil {
		return fmt.Errorf("could not remove matching last answer: %w", err)
	}

	return nil
}

func save(tx *sql.Tx, prompt Prompt, answer int) error {
	stmt, err := tx.Prepare("INSERT INTO answer(element, question, answer, answer_time_ms) values(?,?,?,?)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	t := time.Now()
	elapsed := t.Sub(prompt.startTime)

	if _, err := stmt.Exec(prompt.element.id, prompt.question.id, answer, elapsed.Milliseconds()); err != nil {
		return fmt.Errorf("could not save answer: %w", err)
	}

	return nil
}

func (d Data) Answer(prompt Prompt, answer int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin DB transaction: %w", err)
	}

	defer tx.Rollback()

	// If this question and element match those from the last answer, delete last answer to replace it with this one
	if err := deleteLastIfMatch(tx, prompt); err != nil {
		return fmt.Errorf("could not remove last answer in DB to replace with this one: %w", err)
	}

	// Save the new answer
	if err := save(tx, prompt, answer); err != nil {
		return fmt.Errorf("could not save answer in DB: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction in DB: %w", err)
	}

	return nil
}

func (d Data) Close() error {
	return d.db.Close()
}
