package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/machilan1/go_video/internal/services/auth"
	"github.com/machilan1/go_video/internal/store"
	"github.com/machilan1/go_video/internal/utils/validators"
)

type userLoginForm struct {
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

type userSignupForm struct {
	Name                 string `form:"name"`
	Email                string `form:"email"`
	Password             string `form:"password"`
	validators.Validator `form:"-"`
}

type updateChapterForm struct {
	ChapNum              int    `form:"chapNum"`
	Title                string `form:"title"`
	Description          string `form:"description"`
	FileName             string `form:"fileName"`
	validators.Validator `form:"-"`
}

type createChapterForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	ChapNum     int    `form:"chapNum"`
	CourseID    int    `form:"courseID"`
	FileName    string `form:"fileName"`
}

type createCourseForm struct {
	Title       string `form:"title"`
	Instructor  string `form:"instructor"`
	Tags        []int  `form:"tags"`
	Description string `form:"description"`
}

func (app *application) viewHome(w http.ResponseWriter, r *http.Request) {

	courses, err := app.store.CourseStore.FindMany(store.FindCoursesParams{Page: 1, Limit: 20})
	if err != nil {
		app.serverError(w, r, err)
	}

	courseCards := []courseCardVM{}
	p := &courseCards

	for _, v := range courses {
		*p = append(*p, courseCardVM{
			v,
		})
	}

	res := *p
	data := app.newTemplateData(r)
	data.CourseCards = res
	app.render(w, r, 200, "home.html", data)
}

func (app *application) viewLogin(w http.ResponseWriter, r *http.Request) {
	form := userLoginForm{}
	data := app.newTemplateData(r)

	if data.IsAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	data.Form = form
	app.render(w, r, 200, "login.html", data)
}

func (app *application) viewCourse(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("ID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// find info with id

	c, err := app.store.CourseStore.FindOne(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Course = CourseVM{c}

	app.render(w, r, http.StatusOK, "course.html", data)

}

func (app *application) viewChapter(w http.ResponseWriter, r *http.Request) {
	courseId, err := strconv.Atoi(r.PathValue("courseID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	chapterId, err := strconv.Atoi(r.PathValue("chapterID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// find info with id

	c, err := app.store.CourseStore.FindOne(courseId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var chp store.Chapter

	for _, chapter := range c.Chapters {
		if chapter.ID == chapterId {
			chp = chapter
		}
	}

	data := app.newTemplateData(r)
	data.Course = CourseVM{c}

	data.Chapter = ChapterVM{chp}

	app.render(w, r, http.StatusOK, "course.html", data)

}

func (app *application) viewUserCourse(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	courses, err := app.store.CourseStore.FindMany(store.FindCoursesParams{
		Limit:  20,
		Page:   1,
		UserID: userID})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var courseCards []courseCardVM

	for _, v := range courses {
		courseCards = append(courseCards, courseCardVM{v})
	}

	data := app.newTemplateData(r)
	data.CourseCards = courseCards

	app.render(w, r, http.StatusOK, "user-course.html", data)

}

func (app *application) viewUserCourseDetail(w http.ResponseWriter, r *http.Request) {

	courseID, err := strconv.Atoi(r.PathValue("courseID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	course, err := app.store.CourseStore.FindOne(courseID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Course = CourseVM{course}

	app.render(w, r, http.StatusOK, "user-course-detail.html", data)

}

func (app *application) viewCreateUserCourse(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	tags, err := app.store.TagStore.FindTags()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.Options.Tags = tags

	app.render(w, r, http.StatusOK, "user-course-create.html", data)

}

func (app *application) viewEditUserCourse(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.PathValue("courseID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	course, err := app.store.CourseStore.FindOne(courseID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	tagIDs := []int{}

	for _, v := range course.Tags {
		tagIDs = append(tagIDs, v.ID)
	}

	tags, err := app.store.TagStore.FindTags()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for i := 0; i < len(tags); i++ {
		for _, ID := range tagIDs {
			if ID == tags[i].ID {
				tags[i].Selected = true
			}
		}
	}

	data := app.newTemplateData(r)
	data.Course = CourseVM{course}
	data.Options = OptionsVM{Tags: tags}

	app.render(w, r, http.StatusOK, "user-course-edit.html", data)
}

func (app *application) viewEditCourseChapter(w http.ResponseWriter, r *http.Request) {

	chapterID, err := strconv.Atoi(r.PathValue("chapterID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// get chapter data
	chapter, err := app.store.ChapterStore.FindOneChapter(chapterID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Chapter = ChapterVM{chapter}

	app.render(w, r, http.StatusOK, "user-course-chapter-edit.html", data)
}

func (app *application) viewCreateChapter(w http.ResponseWriter, r *http.Request) {

	courseID, err := strconv.Atoi(r.PathValue("courseID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	data := app.newTemplateData(r)
	data.CourseID = courseID
	app.render(w, r, http.StatusOK, "user-course-chapter-create.html", data)

}

func (app *application) viewAdminUsers(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "admin-users-list.html", data)
}

func (app *application) createCourse(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	tags := []int{}

	for _, v := range r.PostForm["tags"] {
		res, err := strconv.Atoi(v)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		tags = append(tags, res)
	}

	UserID := app.sessionManager.GetInt(r.Context(), "currentUserID")

	form := createCourseForm{
		Title:       r.FormValue("title"),
		Instructor:  r.FormValue("instructor"),
		Tags:        tags,
		Description: r.FormValue("description"),
	}

	// TODO : validate forms

	err = app.store.CourseStore.CreateCourse(store.CreateCourseBody{
		Title:       form.Title,
		Instructor:  form.Instructor,
		Description: form.Description,
		CreatedBy:   UserID,
		Tags:        form.Tags})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "課程新增成功")
	http.Redirect(w, r, "/users/"+fmt.Sprintf("%d", UserID)+"/courses", http.StatusSeeOther)
}

func (app *application) createChapter(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	courseID, err := strconv.Atoi(r.PathValue("courseID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "currentUserID")

	fileName := r.FormValue("fileName")
	description := r.FormValue("description")
	title := r.FormValue("title")
	chapNum, err := strconv.Atoi(r.FormValue("chapNum"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO : 想一下怎麼樣實現 Transaction比較好

	// 新增影片實例
	videoID, err := app.store.VideoStore.CreateVideo(store.CreateVideoBody{
		FileName:    fileName,
		Description: description,
		UpdatedBy:   userID,
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// 新增章節

	_, err = app.store.ChapterStore.Create(courseID, store.CreateChapterBody{
		Title:       title,
		Description: description,
		ChapNum:     chapNum,
		VideoID:     videoID,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "已新增章節")

	http.Redirect(w, r, "/users/"+fmt.Sprintf("%d", userID)+"/courses/"+fmt.Sprintf("%d", courseID), http.StatusSeeOther)

}

func (app *application) editCourse(w http.ResponseWriter, r *http.Request) {

	// TODO : Finish Validation

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	courseID, err := strconv.Atoi(r.PathValue("courseID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	tags := r.PostForm["tags"]
	var tagIds []int

	for _, v := range tags {
		tagID, err := strconv.Atoi(v)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		tagIds = append(tagIds, tagID)
	}

	title := r.PostForm.Get("title")
	instructor := r.PostForm.Get("instructor")
	description := r.PostForm.Get("description")

	err = app.store.CourseStore.UpdateCourse(courseID, store.UpdateCourseBody{
		Title:       title,
		Instructor:  instructor,
		Tags:        tagIds,
		Description: description})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "已更新完成")
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func (app *application) editChapter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	chapNum, err := strconv.Atoi(r.FormValue("chapNum"))

	// TODO : finish validations
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// }

	// form := updateChapterForm{
	// 	ChapNum:     chapNum,
	// 	Title:       r.FormValue("title"),
	// 	Description: r.FormValue("description"),
	// 	FileName:    r.FormValue("fileName"),
	// }

	// form.CheckField(validators.MinChars(form.Title, 1), "title", "Title should be at least 1 letter long")
	// form.CheckField(validators.MaxChars(form.Title, 100), "title", "Title should be at most 100 letters long")
	// form.CheckField(validators.MinChars(form.Description, 1), "description", "Description should be at least 1 letter long")
	// form.CheckField(validators.,"chapNum","")

	userID := r.PathValue("userID")
	courseID := r.PathValue("courseID")
	chapterID, err := strconv.Atoi(r.PathValue("chapterID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	description := r.FormValue("description")
	fileName := r.FormValue("fileName")

	chapter, err := app.store.ChapterStore.FindOneChapter(chapterID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	videoID := chapter.VideoID

	// Update video filename and link
	err = app.store.VideoStore.UpdateVideo(videoID, store.UpdateVideoBody{FileName: fileName})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Update chapter data
	err = app.store.ChapterStore.UpdateChapterInfo(chapterID, store.UpdateChapterBody{Title: title, Description: description, ChapNum: chapNum})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/users/"+userID+"/courses/"+courseID, http.StatusSeeOther)

}

func (app *application) deleteChapter(w http.ResponseWriter, r *http.Request) {
	chapterID, err := strconv.Atoi(r.PathValue("chapterID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	courseID, err := strconv.Atoi(r.PathValue("courseID"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID := app.sessionManager.GetInt(r.Context(), "currentUserID")
	err = app.store.ChapterStore.Delete(chapterID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "已刪除章節")

	http.Redirect(w, r, "/users/"+fmt.Sprintf("%d", userID)+"/courses/"+fmt.Sprintf("%d", courseID), http.StatusSeeOther)

}

func (app *application) deleteCourse(w http.ResponseWriter, r *http.Request) {

	courseID, err := strconv.Atoi(r.PathValue("courseID"))
	userID := r.PathValue("userID")

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.store.CourseStore.Delete(courseID)
	app.sessionManager.Put(r.Context(), "flash", "課程刪除成功")
	http.Redirect(w, r, "/users/"+userID+"/courses", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validators.IsEmail(form.Email), "email", "Invalid email address")
	form.CheckField(validators.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validators.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validators.MinChars(form.Password, 8), "password", "Invalid Password")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	user, err := app.authService.Login(auth.LoginParam{Email: form.Email, Password: form.Password})
	if err != nil {
		data := app.newTemplateData(r)
		data.Form = form
		data.Flash = "帳號或密碼錯誤，請確認後再重試一次"
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "登入成功")
	app.sessionManager.Put(r.Context(), "currentUserID", user.Id)
	app.sessionManager.Put(r.Context(), "canUpdate", user.CanUpdate)
	app.sessionManager.Put(r.Context(), "isAdmin", user.IsAdmin)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) viewSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	if !data.IsAdmin {
		app.clientError(w, http.StatusForbidden)
		return
	}

	data.Form = userSignupForm{}
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data.Flash = flash

	app.render(w, r, http.StatusOK, "signup.html", data)
}

// TODO : make sure you want to keep this handler.
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userSignupForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
		Name:     r.PostForm.Get("name"),
	}

	form.CheckField(validators.NotBlank(form.Email), "title", "This field cannot be blank")
	form.CheckField(validators.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validators.NotBlank(form.Name), "name", "This field cannot be blank")

	form.CheckField(validators.MinChars(form.Password, 8), "password", "Password should be larger than 8 charactors")
	form.CheckField(validators.MaxChars(form.Password, 50), "password", "Password should be less than 50 charactors")
	form.CheckField(validators.CharAndNumOnly(form.Password), "password", "Password should only contain letters and numbers")
	form.CheckField(validators.IsEmail(form.Email), "email", "Not a valid email address")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.authService.SignUp(auth.SignUpParam{Password: form.Password, Name: form.Name, Email: form.Email})

	if err != nil {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusBadRequest, "signup.html", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "註冊成功 請登入")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *application) upload(w http.ResponseWriter, r *http.Request) {

	f, h, err := r.FormFile("file")
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = os.WriteFile("statics/videos/"+h.Filename, bs, 0644)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	p := struct {
		FileName string `json:"fileName"`
	}{
		FileName: h.Filename,
	}

	js, err := json.Marshal(p)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = w.Write(js)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {

	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "currentUserID")
	app.sessionManager.Put(r.Context(), "flash", "您已登出")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
