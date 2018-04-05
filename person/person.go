package person

import "strconv"

func (m Person) TestString() string {
	result := m.Name.Family + m.Name.Personal

	for i, e := range m.Email {
		result = result + strconv.Itoa(i) + e.Address + e.Kind
	}
	return result
}
