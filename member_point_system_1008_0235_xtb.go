// 代码生成时间: 2025-10-08 02:35:26
package main

import (
    "encoding/json"
    "fmt"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
)

// MemberPoint is the structure for member points
type MemberPoint struct {
    Id       int    `orm:"auto"`
    MemberId int    `orm:"index"`
    Points   int    `orm:"default(0)"`
}

// MemberPointService handles business logic for member points
type MemberPointService struct {
    o    orm.Ormer
    sync bool
}

// NewMemberPointService returns a new instance of MemberPointService
func NewMemberPointService() *MemberPointService {
    return &MemberPointService{
        sync: true,
    }
}

// AddPoint adds points to a member
func (s *MemberPointService) AddPoint(memberId int, points int) error {
    if points < 0 {
        return fmt.Errorf("points cannot be negative")
    }
    s.o.LoadRelated(&MemberPoint{MemberId: memberId}, "Points")
    if _, err := s.o.Update(&MemberPoint{MemberId: memberId, Points: MemberPoint(s.o).Points + points}, "Points"); err != nil {
        return err
    }
    return nil
}

// GetPoints retrieves points for a member
func (s *MemberPointService) GetPoints(memberId int) (int, error) {
    var points int
    if err := s.o.QueryTable(&MemberPoint{}).Filter("MemberId", memberId).One("", &points); err != nil {
        return 0, err
    }
    return points, nil
}

func main() {
    beego.TestBeegoInit("./")
    o := orm.NewOrm()
    defer o.Db.Close()
    service := NewMemberPointService()
    // Register model
    err := o.SyncDB()
    if err != nil {
        fmt.Println(err)
        return
    }
    // Testing adding points
    memberId := 1
    points := 100
    if err := service.AddPoint(memberId, points); err != nil {
        fmt.Println("Error adding points: ", err)
    } else {
        fmt.Printf("Added %d points to member %d
", points, memberId)
    }
    // Testing getting points
    retrievedPoints, err := service.GetPoints(memberId)
    if err != nil {
        fmt.Println("Error getting points: ", err)
    } else {
        fmt.Printf("Member %d has %d points
", memberId, retrievedPoints)
    }
}
