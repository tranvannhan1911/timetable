package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)



type Busy struct {
	num_busy     int
	lession_busy []Lession
}

type Teacher struct {
	id string
	busy Busy
}

type Class struct {
	id   int
	name string
}

type Subject struct {
	id string
}

type Session struct {
	id          string
	session     string
	day_of_week int
}

type Lession struct {
	id               int
	name             string
	lessionOfSession int
	session          Session
}

type Assignment struct {
	id              string
	class           Class
	subject         Subject
	teacher         Teacher
	numberOfLession int
	length          int
}

type TimeTableAssignment struct {
	teacher Teacher
	subject Subject
}

type TimeTable struct {
	timeTable [][]TimeTableAssignment
	fitness   int
}

var classes []Class
var lessions []Lession
var teachers []Teacher

func (tt TimeTable) weaker(tt2 TimeTable) bool {
	return tt.fitness < tt2.fitness
}

type ByFitness []TimeTable

func (s ByFitness) Len() int {
	return len(s)
}

func (s ByFitness) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByFitness) Less(i, j int) bool {
	return s[i].weaker(s[j])
}

func initLession(sessions *[]Session, lessions *[]Lession) {
	lession_idx := 0
	for i := 2; i <= 6; i++ {
		am := Session{
			id:          "T" + strconv.Itoa(i) + " - AM",
			session:     "AM",
			day_of_week: i,
		}
		*sessions = append(*sessions, am)
		for j := 1; j <= 5; j++ {
			lession := Lession{
				id: lession_idx,
				// name: ,
				lessionOfSession: j,
				session:          am,
			}
			*lessions = append(*lessions, lession)
			lession_idx++
		}

		pm := Session{
			id:          "T" + strconv.Itoa(i) + " - PM",
			session:     "PM",
			day_of_week: i,
		}
		*sessions = append(*sessions, pm)
		for j := 1; j <= 5; j++ {
			lession := Lession{
				id: lession_idx,
				// name: ,
				lessionOfSession: j,
				session:          pm,
			}
			*lessions = append(*lessions, lession)
			lession_idx++
		}
	}
}

func input(assignments *[]Assignment, classes *[]Class, teachers *[]Teacher, subjects *[]Subject) {
	file, err := os.Open("PC_HK1.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)

	classIdxCounter := 0
	classIndexes := make(map[string]int)
	teacherCheck := make(map[string]bool)
	subjectCheck := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != 6 {
			fmt.Println("Invalid line:", line)
			continue
		}
		class := Class{
			id:   0,
			name: fields[1],
		}

		if _, ok := classIndexes[class.name]; ok {
			class.id = classIndexes[class.name]
		} else {
			class.id = classIdxCounter
			*classes = append(*classes, class)
			classIndexes[fields[1]] = classIdxCounter
			classIdxCounter++
		}
		subject := Subject{
			id: fields[2],
		}
		teacher := Teacher{
			id: fields[3],
		}
		var numberOfLession, _ = strconv.Atoi(fields[4])
		var length, _ = strconv.Atoi(fields[5])
		var an_assignment = Assignment{
			id:              fields[0],
			class:           class,
			subject:         subject,
			teacher:         teacher,
			numberOfLession: numberOfLession,
			length:          length,
		}
		*assignments = append(*assignments, an_assignment)

		if !teacherCheck[teacher.id] {
			*teachers = append(*teachers, teacher)
		}
		if !subjectCheck[subject.id] {
			*subjects = append(*subjects, subject)
		}
		teacherCheck[teacher.id] = true
		subjectCheck[subject.id] = true
	}
}


func inputTeacherBusy(teacher *[]Teacher) {
	// 	// col structure of the file PC: Mã - Mã GV - Thứ - Buổi - Tiết
	file, err := os.Open("GV_Busy.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	teachers_busy := make(map[string][]Lession)

	lst_teacher_busy := []string{}

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != 5 {
			fmt.Println("Invalid line:", line)
			continue
		}
		day_of_week, _ := strconv.Atoi(fields[2])
		lessionOfSession, _ := strconv.Atoi(fields[4])
		var lession_busy = Lession{
			lessionOfSession: lessionOfSession,
			session: Session{
				day_of_week: day_of_week,
				session:     fields[3],
			},
		}
		var _, ok = teachers_busy[fields[1]]
		if !ok {
			lst_teacher_busy = append(lst_teacher_busy, fields[1])
		}
		teachers_busy[fields[1]] = append(teachers_busy[fields[1]], lession_busy)
	}
	teacher_busy_check := make(map[string]bool)

	var tmp_teachers []Teacher

	for _, value := range *teacher {
		val, ok := teachers_busy[value.id]
		if ok {
			busy := Busy{
				num_busy:     len(val),
				lession_busy: val,
			}
			// value.busy.num_busy = len(val)
			// fmt.Println("test:", val)
			// for _, val1 := range val {
			// 	value.busy.session_busy = append(value.busy.session_busy, val1)
			// }
			value.busy = busy
			teacher_busy_check[value.id] = true
		}
		tmp_teachers = append(tmp_teachers, value)
	}
	*teacher = tmp_teachers
	for _, value := range lst_teacher_busy {
		_, ok := teacher_busy_check[value]
		if !ok {
			fmt.Printf("Teacher ID Busy does not exist: %s\n", value)
		}
	}
}

func printAssignments(assignments []Assignment) {
	for index, value := range assignments {
		fmt.Printf("Index: %d, Value: %s\n", index, value)
	}
	fmt.Println()
}

func printLession(lessions []Lession) {
	for _, value := range lessions {
		fmt.Printf("id: %d, Session: %s, LOS: %d\n", value.id, value.session.id, value.lessionOfSession)
	}
	fmt.Println()
}

func printClass(classes []Class) {
	for _, value := range classes {
		fmt.Printf("id: %d, name: %d\n", value.id, value.name)
	}
	fmt.Println()
}

func printSubject(subjects []Subject) {
	for index, value := range subjects {
		fmt.Printf("Index: %d, Value: %d\n", index, value.id)
	}
	fmt.Println()
}

func printTeacher(teachers []Teacher) {
	for index, value := range teachers {
		fmt.Printf("Index: %d, Value: %d\n", index, value.id)
	}
	fmt.Println()
}

func (tt *TimeTable) writeToCSV(classes []Class, lessions []Lession, filename string) error {
	// Mở file để ghi
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Tạo một writer CSV từ file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var classHeaders []string
	classHeaders = append(classHeaders, "Class")
	for _, class := range classes {
		classHeaders = append(classHeaders, class.name)
	}
	writer.Write(classHeaders)

	// Duyệt qua từng hàng của bảng thời gian và ghi vào file CSV
	for index, row := range tt.timeTable {
		var record []string
		record = append(record, fmt.Sprintf("%s - Tiet %d", lessions[index].session.id, lessions[index].lessionOfSession))
		for _, assignment := range row {
			if assignment.teacher.id != "" && assignment.subject.id != "" {
				record = append(record, fmt.Sprintf("%s - %s", assignment.teacher.id, assignment.subject.id))
			} else {
				record = append(record, "")
			}
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func initialEmptyTimeTable(classes []Class, lessions []Lession) TimeTable {
	timetable := TimeTable{
		timeTable: make([][]TimeTableAssignment, len(lessions)), // Khởi tạo slice con đầu tiên với độ dài 3
	}

	// Khởi tạo các slice con với độ dài 4 và giá trị mặc định là TimeTableAssignment{}
	for i := range timetable.timeTable {
		timetable.timeTable[i] = make([]TimeTableAssignment, len(classes))
	}
	return timetable
}

func (a TimeTableAssignment) isEmpty() bool {
	if a.teacher.id == "" || a.subject.id == "" {
		return true
	}
	return false
}

func (tt TimeTable) findFirstEmptyRow(col int) int {
	for row := 0; row < len(tt.timeTable); row++ {
		// fmt.Println(timeTable.timeTable[row][col].teacher.id)
		if tt.timeTable[row][col].isEmpty() {
			return row
		}
	}
	return -1 // Trả về -1 nếu không tìm thấy dòng rỗng
}

func (tt *TimeTable) addTimeTableAssignment(class Class, lession Lession, teacher Teacher, subject Subject) {
	col := class.id
	row := lession.id
	tt.timeTable[row][col] = TimeTableAssignment{
		teacher: teacher,
		subject: subject,
	}
}

func (tt TimeTable) getRandomLessionAssignment(lessions []Lession, class Class) Lession {
	// TODO: Tối ưu hàm này bằng cách lấy các tiết dạy đang trống rồi mới random các tiết đó
	counter := 0
	for counter <= 2*len(lessions) {
		row := rand.Intn(len(lessions))
		if tt.timeTable[row][class.id].isEmpty() {
			return lessions[row]
		}
		counter++
	}
	return Lession{}
}

func initTimeTable(assignments []Assignment, classes []Class, teachers []Teacher, subjects []Subject, lessions []Lession, timeTable *TimeTable) {
	for _, an_assignment := range assignments {
		for i := 0; i < an_assignment.numberOfLession; i++ {
			lession := timeTable.getRandomLessionAssignment(lessions, an_assignment.class)
			timeTable.addTimeTableAssignment(an_assignment.class, lession, an_assignment.teacher, an_assignment.subject)
		}
	}
	timeTable.calFitness()
}

// check 1 tiết cụ thể, xem 1 giáo viên có dạy trùng hay không, return số lần bị trùng
func (tt TimeTable) checkTrungTiet(teacher Teacher, lession Lession) int {
	row := lession.id
	counter := 0
	for col := 0; col < len(tt.timeTable[0]); col++ {
		if tt.timeTable[row][col].teacher.id == teacher.id {
			counter++
		}
	}
	return counter - 1
}

// số lượng giáo viên bị trùng tiết
func (tt TimeTable) calTrungTiet() int {
	fitness := 0
	for row := 0; row < len(tt.timeTable); row++ {
		teacherCheck := make(map[string]bool)
		for col := 0; col < len(tt.timeTable[row]); col++ {
			if tt.timeTable[row][col].isEmpty() {
				continue
			}
			if teacherCheck[tt.timeTable[row][col].teacher.id] {
				fitness++
			}

			teacherCheck[tt.timeTable[row][col].teacher.id] = true
		}
		// fmt.Println(fitness)
	}
	// fmt.Println("fitness", fitness)
	return fitness
}

// kiểm tra tiết lủng
func (tt TimeTable) calTietLung() int {
	fitness := 0
	for col := 0; col < len(tt.timeTable[0]); col++ {
		for row := 0; row < len(tt.timeTable); row++ {
			if tt.timeTable[row][col].isEmpty() {
				var c1, c2 bool
				k_truoc := row
				for k_truoc >= (row/5)*5 {
					if !tt.timeTable[k_truoc][col].isEmpty() { // Tiết trước đó có tiết
						c1 = true
						break
					}
					k_truoc--
				}
				k_sau := row
				for k_sau < (row/5+1)*5 {
					if !tt.timeTable[k_sau][col].isEmpty() { // Tiết sau đó có tiết
						c2 = true
						break
					}
					k_sau++
				}

				if c1 && c2 {
					fitness++
				}
			}
		}
	}
	return fitness
}

// kiểm tra giáo viên chỉ dạy 1 tiết trong 1 buổi
func (tt TimeTable) calBuoiDay1Tiet() int {
	fitness := 0
	teacherCheck := make(map[string]int)

	for row := 0; row < len(tt.timeTable); row++ {
		for col := 0; col < len(tt.timeTable[row]); col++ {
			if !tt.timeTable[row][col].isEmpty() {
				teacherCheck[tt.timeTable[row][col].teacher.id]++
			}
		}
		if (row+1)%5 == 0 {
			for _, val := range teacherCheck {
				if val == 1 {
					fitness++
				}
			}
			teacherCheck = make(map[string]int) // Reset the map
		}
	}
	return fitness
}

// kiểm tra số tiết tối thiểu của 1 lớp
func (tt TimeTable) calTietToiThieu() int {
	fitness := 0
	for col := 0; col < len(tt.timeTable[0]); col++ {
		t := 0
		for row := 0; row < len(tt.timeTable); row++ {
			if !tt.timeTable[row][col].isEmpty() {
				t++
			}
			if (row+1)%5 == 0 {
				if t < 2 && t > 0 {
					fitness++
				}
				t = 0
			}
		}
	}
	return fitness
}

// kiểm tra tối đa môn

func (tt TimeTable) calToiDaMon() int {
	fitness := 0
	for col := 0; col < len(tt.timeTable[0]); col++ {
		t := 0
		subjectCheck := make(map[string]bool)
		for row := 0; row < len(tt.timeTable); row++ {
			if !tt.timeTable[row][col].isEmpty() {
				if !subjectCheck[tt.timeTable[row][col].subject.id] {
					t++
				}
				subjectCheck[tt.timeTable[row][col].subject.id] = true
			}
			if (row+1)%5 == 0 {
				if t > 4 {
					fitness++
				}
				t = 0
				subjectCheck = make(map[string]bool) // Reset the map
			}
		}
	}
	return fitness
}


var teacher_lession_busy = make(map[string]map[string]bool)

func createMapLessionTeacherBusy() {
	for _, value := range teachers {
		for _, value1 := range value.busy.lession_busy {
			s := "T" + strconv.Itoa(value1.session.day_of_week) + " - " + value1.session.session + strconv.Itoa(value1.lessionOfSession)
			_, ok := teacher_lession_busy[s]
			if !ok {
				teacher_lession_busy[s] = make(map[string]bool)
			}
			teacher_lession_busy[s][value.id] = true
		}
	}
}

func (tt TimeTable) calTeacherBusy() int {
	fitness := 0

	for index, row := range tt.timeTable {
		for _, assignment := range row {
			if assignment.teacher.id != "" && assignment.subject.id != "" {
				// record = append(record, fmt.Sprintf("%s - %s", assignment.teacher.id, assignment.subject.id))
				s := "T" + strconv.Itoa(lessions[index].session.day_of_week) + " - " + lessions[index].session.session + strconv.Itoa(lessions[index].lessionOfSession)
				if teacher_lession_busy[s][assignment.teacher.id] {
					fitness += 1
				}
			}
		}
	}

	// for row := 0; row < len(tt.timeTable); row++ {
	// 	teacherCheck := make(map[string]bool)
	// 	for col := 0; col < len(tt.timeTable[row]); col++ {
	// 		if tt.timeTable[row][col].isEmpty() {
	// 			continue
	// 		}
	// 		if teacherCheck[tt.timeTable[row][col].teacher.id] {
	// 			fitness++
	// 		}

	// 		teacherCheck[tt.timeTable[row][col].teacher.id] = true
	// 	}
	// 	// fmt.Println(fitness)
	// }
	// // fmt.Println("fitness", fitness)
	return fitness
}

// đánh giá sự tối ưu của 1 TKB
func (tt *TimeTable) calFitness() {
	// hiện tại chỉ check trùng tiết hay không
	var res int = 0
	res += tt.calTrungTiet() * 999
	res += tt.calTietLung() * 600
	res += tt.calTeacherBusy() * 600
	res += tt.calBuoiDay1Tiet() * 10
	res += tt.calTietToiThieu() * 20
	res += tt.calToiDaMon() * 10
	tt.fitness = res
}

// đánh giá tối ưu của 1 tiết của 1 lớp trong thời khóa biểu
func (tt TimeTable) calAssignmentFitness(class Class, lession Lession) int {
	// hiện tại chỉ check xem gv đó có bị trùng tiết ở lớp khác hay không
	return tt.checkTrungTiet(tt.timeTable[class.id][lession.id].teacher, lession)
}

func (tt *TimeTable) swapAssignment(class Class, lessionA Lession, lessionB Lession) {
	tmp := tt.timeTable[lessionA.id][class.id]
	tt.timeTable[lessionA.id][class.id] = tt.timeTable[lessionB.id][class.id]
	tt.timeTable[lessionB.id][class.id] = tmp
}

func (tt TimeTable) getDuplicateLessionsOfClass(class Class) []Lession {
	var dupLessions []Lession
	for row := 0; row < len(tt.timeTable); row++ {
		if !tt.timeTable[row][class.id].isEmpty() {
			dup := tt.checkTrungTiet(tt.timeTable[row][class.id].teacher, lessions[row])
			if dup > 0 {
				dupLessions = append(dupLessions, lessions[row])
			}
		}
	}
	return dupLessions
}

func (tt TimeTable) getTeacherBusyOfClass(class Class) []Lession {
	var busyLession []Lession

	for row := 0; row < len(tt.timeTable); row++ {
		if !tt.timeTable[row][class.id].isEmpty() {
			s := "T" + strconv.Itoa(lessions[row].session.day_of_week) + " - " + lessions[row].session.session + strconv.Itoa(lessions[row].lessionOfSession)
			if teacher_lession_busy[s][tt.timeTable[row][class.id].teacher.id] {
				busyLession = append(busyLession, lessions[row])
			}
		}
	}

	return busyLession
}

func (tt TimeTable) clone() TimeTable {
	newTimeTable := initialEmptyTimeTable(classes, lessions)
	for row := 0; row < len(tt.timeTable); row++ {
		copy(newTimeTable.timeTable[row], tt.timeTable[row])
	}
	return newTimeTable
}

func (tt TimeTable) improve() TimeTable {
	// clone ra 1 timetable mới
	newTimeTable := tt.clone()
	newTimeTable.calFitness()
	// chọn ngẫu nhiên 1 lớp trong TKB để thực hiện cải thiện
	col := rand.Intn(len(tt.timeTable[0]))
	p := rand.Intn(len(tt.timeTable))
	var idx1, idx2 int
	dupLessions := newTimeTable.getDuplicateLessionsOfClass(classes[col])

	busyLession := newTimeTable.getTeacherBusyOfClass(classes[col])

	// vailTeacherBusy :=
	// nếu trong lớp đó không có tiết nào bị trùng thì skip, có thể không phù hợp với bài toán nhiều ràng buộc
	if len(dupLessions) == 0 && len(busyLession) == 0 {
		return newTimeTable
	}

	if len(dupLessions) > 0 {
		if p <= 30 {
			idx1 = rand.Intn(len(tt.timeTable))
			idx2 = rand.Intn(len(tt.timeTable))
		} else { // xác suất 70% là thực hiện swap 1 tiết bị trùng và 1 tiết chọn ngẫu nhiên
			idx1 = dupLessions[rand.Intn(len(dupLessions))].id
			idx2 = rand.Intn(len(tt.timeTable))
		}
		// thực hiện swap 2 tiết trong 1 lớp
		newTimeTable.swapAssignment(classes[col], lessions[idx1], lessions[idx2])
	}

	if len(busyLession) > 0 {
		if p <= 30 {
			idx1 = rand.Intn(len(tt.timeTable))
			idx2 = rand.Intn(len(tt.timeTable))
		} else { // xác suất 70% là thực hiện swap 1 tiết bị trùng và 1 tiết chọn ngẫu nhiên
			idx1 = busyLession[rand.Intn(len(busyLession))].id
			idx2 = rand.Intn(len(tt.timeTable))
		}
		// thực hiện swap 2 tiết trong 1 lớp
		newTimeTable.swapAssignment(classes[col], lessions[idx1], lessions[idx2])
	}

	newTimeTable.calFitness()
	return newTimeTable
}

func geneticAlgo(assignments []Assignment, classes []Class, teachers []Teacher, subjects []Subject, lessions []Lession) {
	const POPULATION_SIZE = 100
	generation := 0

	var population []TimeTable
	for i := 0; i < POPULATION_SIZE; i++ {
		timeTable := initialEmptyTimeTable(classes, lessions)
		initTimeTable(assignments, classes, teachers, subjects, lessions, &timeTable)
		timeTable.calFitness()
		population = append(population, timeTable)
	}

	for generation <= 1000 {
		sort.Sort(ByFitness(population))
		// for i := 0; i < 10; i++ {
		// 	fmt.Println(population[i].fitness)
		// }
		// fmt.Println()
		if population[0].fitness <= 0 {
			break
		}
		var newGeneration []TimeTable
		s10 := (10 * POPULATION_SIZE) / 100
		for i := 0; i < s10; i++ {
			newGeneration = append(newGeneration, population[i])
		}
		s90 := (90 * POPULATION_SIZE) / 100
		s50 := (20 * len(population)) / 100
		for i := 0; i < s90; i++ {
			r1 := rand.Intn(s50)
			parent1 := population[r1]
			offspring := parent1.improve()
			newGeneration = append(newGeneration, offspring)
		}
		// population := newGeneration
		copy(population, newGeneration)
		fmt.Println(generation, population[0].fitness)
		generation++
	}
	population[0].writeToCSV(classes, lessions, "aoutput.csv")
	fmt.Println(generation, population[0].fitness)
	for col := 0; col < len(population[0].timeTable[0]); col++ {
		dup := population[0].getDuplicateLessionsOfClass(classes[col])
		fmt.Println(dup)
	}
}

// countBuoiDay1Tiet
func main() {
	rand.Seed(time.Now().UnixNano())
	// Create a slice to store PhanCong structs
	assignments := []Assignment{}
	classes = []Class{}
	lessions = []Lession{}

	teachers := []Teacher{}
	subjects := []Subject{}
	sessions := []Session{}

	initLession(&sessions, &lessions)
	// printLession(lessions)

	input(&assignments, &classes, &teachers, &subjects)
	inputTeacherBusy(&teachers)

	geneticAlgo(assignments, classes, teachers, subjects, lessions)
	// printAssignments(assignments)
	// extractInfo(assignments, &classes, &teachers, &subjects)

	// printTeacher(teachers)
	// printClass(classes)
	// printSubject(subjects)

	// timeTable := initialEmptyTimeTable(classes, lessions)
	// initTimeTable(assignments, classes, teachers, subjects, lessions, &timeTable)
	// timeTable.calFitness()
	// fmt.Println("fitness: ", timeTable.fitness)

	// newTimeTable := timeTable.mate()
	// fmt.Println("fitness: ", newTimeTable.fitness)

	// fmt.Println(timeTable.timeTable)

	// timeTable.writeToCSV(classes, lessions, "output.csv")

}