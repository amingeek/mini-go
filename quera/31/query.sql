use quera;

show tables;

select * from students;

select * from student_courses;

SELECT
  *
FROM students
WHERE NOT EXISTS (SELECT
  1
FROM student_courses
WHERE course_id = 7
AND student_id = students.id);


select students.id, students.name
   from students left join student_courses on students.id = student_courses.student_id
    where not exists(
        SELECT
  1
FROM student_courses
WHERE course_id = 7
AND student_id = students.id
    );



SELECT students.id, students.name
FROM students
LEFT JOIN student_courses
    ON students.id = student_courses.student_id
    AND student_courses.course_id = 7
WHERE student_courses.student_id IS NULL
ORDER BY students.id;

