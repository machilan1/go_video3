
insert into app_user (name, email, password,is_admin) values ('admin', 'admin@gmail.com', '$2a$12$NjDc43MPK8GWjVo2kVNRUuwT.f/JeFtnHvryGoRxyF6xSip8wvTDq',true);
insert into app_user (name, email, password) values ('student','student@gmail.com','$2a$12$NjDc43MPK8GWjVo2kVNRUuwT.f/JeFtnHvryGoRxyF6xSip8wvTDq');
insert into app_user (name, email, password, can_update) values ('teacher','teacher@gmail.com','$2a$12$NjDc43MPK8GWjVo2kVNRUuwT.f/JeFtnHvryGoRxyF6xSip8wvTDq',true);

insert into course (name, description, instructor_name, created_by) values ('Math 101', 'Math course', 'Bill Gates',3);
insert into course (name, description, instructor_name, created_by) values ('Algorithm 101', 'Algorithm course', 'Lebron James',3);
insert into course (name, description, instructor_name, created_by) values ('Backend 101', 'Backend course', 'James Schwantz',3);

insert into tag (label, created_by) values ('math',3);
insert into tag (label, created_by) values ('algorithm',3);
insert into tag (label, created_by) values ('backend',3);
insert into tag (label, created_by) values ('golang',3);
insert into tag (label, created_by) values ('中文',3);

insert into course_tag (course_id, tag_id) values (1,1);
insert into course_tag (course_id, tag_id) values (2,2);
insert into course_tag (course_id, tag_id) values (3,3);
insert into course_tag (course_id, tag_id) values (3,4);
insert into course_tag (course_id, tag_id) values (3,5);

insert into video ( file_name,updated_by) values ('1.mkv',3);
insert into video ( file_name,updated_by) values ('2.mkv',3);

insert into chapter (title,chap_num,description,course_id,video_id) values ('Introduction',1,'Introduction to Math',1,1);
insert into chapter (title,chap_num,description,course_id,video_id) values ('Linear Algebra',2,'Introduction to Linear Algebra' ,1,2);