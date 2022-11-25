drop database if exists dbTodoList;
create database if not exists dbtodolist;
use dbtodolist;

CREATE TABLE tbTodos (
    id INT PRIMARY KEY AUTO_INCREMENT,
    todoName CHAR(64) NOT NULL,
    todoDescription TEXT,
    isChecked bool not null
);

-- INSERT EXEMPLE
/*
insert into tbTodos
values
(default, "Database", "Create a Mysql database to work with the api and frontend", false);
*/
