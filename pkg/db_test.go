package pkg_test

//var db = pkg.NewDatabase("kuto", "localhost:3306", "root", "root")

//func TestSelect(t *testing.T) {
//	var category models.ArticleCategory
//	err := db.SelectByID(&category, 1)
//	t.Log("category=", category)
//	if err != nil {
//		t.Error("select by id failed, err=", err)
//	}
//
//	var categories []models.ArticleCategory
//	err = db.Select(&categories, "id>?", 5)
//	t.Log("categories=", categories)
//	if err != nil {
//		t.Error("select failed, err=", err)
//	}
//}
//
//func TestInsert(t *testing.T) {
//	c1 := &models.ArticleCategory{
//		Name:     "t1",
//		Desc:     "d1",
//		Keywords: "k1",
//	}
//	c2 := &models.ArticleCategory{
//		Name:     "t2",
//		Desc:     "d2",
//		Keywords: "k2",
//	}
//
//	err := db.Insert(c1, c2)
//	if err != nil {
//		t.Error("insert failed err=", err)
//	}
//}
//
//func TestUpdate(t *testing.T) {
//	c1 := &models.ArticleCategory{
//		Name:     "t1",
//		Desc:     "d1",
//		Keywords: "k1",
//	}
//	err := db.Insert(c1)
//	if err != nil {
//		t.Error("update insert failed, err=", err)
//		return
//	}
//	err = db.SelectByID(c1, c1.ID)
//	if err != nil {
//		t.Error("select failed, err=", err)
//		return
//	}
//	c1.Name = "t2"
//	err = db.Update(c1, "name")
//	if err != nil {
//		t.Error("update failed, err=", err)
//		return
//	}
//	db.SelectByID(c1, c1.ID)
//	if c1.Name != "t2" {
//		t.Log("update failed, not equal")
//	}
//}
//
//func TestDelete(t *testing.T) {
//	c1 := &models.ArticleCategory{
//		Name:     "t1",
//		Desc:     "d1",
//		Keywords: "k1",
//	}
//	err := db.Insert(c1)
//	if err != nil {
//		t.Error("delete insert failed, err=", err)
//		return
//	}
//	err = db.Delete(c1)
//	if err != nil {
//		t.Error("delete failed, err=", err)
//	}
//}
//
//func TestLog(t *testing.T) {
//	logger := pkg.NewLogger(os.Stdout, true)
//	err := errors.New("err test")
//	logger.E("test=%s", err)
//}
