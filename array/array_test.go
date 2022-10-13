package array

import (
	"math/rand"
	"testing"
)

type Player struct {
	id int
}

func TestNew(t *testing.T) {
	list := New[Player](10)
	if len(list.nodes) != 10 {
		t.Fatal("list size error")
	}
}

func TestAdd(t *testing.T) {
	list := New[Player](10)
	for i := 0; i < 10; i++ {
		p := Player{
			id: i,
		}
		index, ok := list.Add(&p)
		if !ok {
			t.Fatal("add failed")
		}
		if index != i {
			t.Fatal("add index wrong", index, i)
		}
	}

	p := Player{
		id: 10,
	}
	_, ok := list.Add(&p)
	if ok {
		t.Fatal("add out of size")
	}

	pid := make(map[int]bool)
	list.ForRange(func(p *Player) {
		pid[p.id] = true
	})
	for i := 0; i < 10; i++ {
		if !pid[i] {
			t.Fatal("pid no exist", i)
		}
	}
	if len(pid) != 10 {
		t.Fatal("pid wrong", pid)
	}
}

func TestDel(t *testing.T) {
	list := New[Player](10)
	pmap := make(map[int]Player)
	for i := 0; i < 10; i++ {
		p := Player{
			id: i,
		}
		index, _ := list.Add(&p)
		pmap[index] = p
	}

	list.Del(0)
	list.Del(5)
	list.Del(9)

	pid := make(map[int]bool)
	list.ForRange(func(p *Player) {
		if pmap[0].id == p.id || pmap[5].id == p.id || pmap[9].id == p.id {
			t.Fatal("not del", p.id)
		}
		pid[p.id] = true
	})

	if len(pid) != 7 {
		t.Fatal("pid wrong", pid)
	}
}

var players []Player

func init() {
	players = make([]Player, 10)
	for i := range players {
		players[i].id = i
	}
}

func BenchmarkRandomAddDel(b *testing.B) {
	list := New[Player](100)

	for i := 0; i < b.N; i++ {
		if list.length < 100 {
			var p Player
			p.id = i
			list.Add(&p)
		} else {
			list.Del(rand.Intn(100))
		}
	}
}

func BenchmarkTestDelAdd(b *testing.B) {
	list := New[Player](10)
	pmap := make(map[int]*Player)
	indexMap := make(map[*Player]int)
	for i := 0; i < 10; i++ {
		players[i].id = i
		index, _ := list.Add(&players[i])
		pmap[index] = &players[i]
		indexMap[&players[i]] = index
	}

	p0, p5, p9 := &players[0], &players[5], &players[9]

	for i := 0; i < b.N; i++ {
		list.Del(indexMap[p0])
		list.Del(indexMap[p5])
		list.Del(indexMap[p9])

		index, _ := list.Add(p0)
		pmap[index] = p0
		indexMap[p0] = index
		index, _ = list.Add(p5)
		pmap[index] = p5
		indexMap[p5] = index
		index, _ = list.Add(p9)
		pmap[index] = p9
		indexMap[p9] = index
	}

	pid := make(map[int]bool)
	list.ForRange(func(p *Player) {
		pid[p.id] = true
	})
	for i := 0; i < 10; i++ {
		if !pid[i] {
			b.Fatal("pid no exist", i)
		}
	}
	if len(pid) != 10 {
		b.Fatal("pid wrong", pid)
	}
}
