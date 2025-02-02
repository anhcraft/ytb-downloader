package request

// No locking required on RequestTable
// As long as the caller is the UI thread

type Table struct {
	requests []*Request
}

var tableInstance = &Table{}

func GetTable() *Table {
	return tableInstance
}

func (rt *Table) Count() int {
	return len(rt.requests)
}

func (rt *Table) AddBulk(req []*Request) {
	rt.requests = append(rt.requests, req...)
}

func (rt *Table) Remove(req *Request) {
	for i, r := range rt.requests {
		if r == req {
			rt.requests = append(rt.requests[:i], rt.requests[i+1:]...)
			break
		}
	}
}

func (rt *Table) Get(index int) *Request {
	if index >= 0 && index < len(rt.requests) {
		return rt.requests[index]
	}
	return nil
}

func (rt *Table) Clear() {
	rt.requests = rt.GetAllByStatus(StatusDownloading)
}

func (rt *Table) GetAllByStatus(status uint32) []*Request {
	var result []*Request
	for _, r := range rt.requests {
		if r.Status() == status {
			result = append(result, r)
		}
	}
	return result
}
