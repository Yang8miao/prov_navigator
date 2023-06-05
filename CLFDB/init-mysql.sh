#!/bin/bash
mysql -u root -proot_clf#2023 clf_db < "/docker-entrypoint-initdb.d/init_tables.sql"