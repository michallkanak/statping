package main

import "time"

type Hit struct {
	Id        int
	Metric    int
	Value     float64
	CreatedAt time.Time
}

func (s *Service) Hits() []Hit {
	var tks []Hit
	rows, err := db.Query("SELECT * FROM hits WHERE service=$1 ORDER BY id DESC LIMIT 256", s.Id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var tk Hit
		err = rows.Scan(&tk.Id, &tk.Metric, &tk.Value, &tk.CreatedAt)
		if err != nil {
			panic(err)
		}
		tks = append(tks, tk)
	}
	return tks
}

func (s *Service) SelectHitsGroupBy(group string) []Hit {
	var tks []Hit
	rows, err := db.Query("SELECT date_trunc('$1', created_at), -- or hour, day, week, month, year count(1) FROM hits WHERE service=$2 group by 1", group, s.Id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var tk Hit
		err = rows.Scan(&tk.Id, &tk.Metric, &tk.Value, &tk.CreatedAt)
		if err != nil {
			panic(err)
		}
		tks = append(tks, tk)
	}
	return tks
}

func (s *Service) TotalHits() int {
	var amount int
	db.QueryRow("SELECT COUNT(id) FROM hits WHERE service=$1;", s.Id).Scan(&amount)
	return amount
}

func (s *Service) Sum() float64 {
	var amount float64
	db.QueryRow("SELECT SUM(latency) FROM hits WHERE service=$1;", s.Id).Scan(&amount)
	return amount
}
