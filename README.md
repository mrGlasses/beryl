# Beryl SQL Helper
![BerylSQL](https://user-images.githubusercontent.com/15248665/201835936-e13f65ff-c267-4569-824d-e30d09769490.png)


Do you want a simple SQL file executer? This is what we have. (Inspired by RedGate's [Flyway](https://flywaydb.org/))

We do not create ghost tables like [gh-ost](https://github.com/github/gh-ost) or create versions of tables like [Skeema](https://www.skeema.io/)/[SkeeFree](https://github.com/github/gh-mysql-tools/tree/master/skeefree), we simply run scripts.

Keep it simple, keep it pretty.

## Installation
1. Obviously install Go if you don't
2. ~~Because of damned Oracle driver~~ Mainly you have to install GCC([Windows](https://code.visualstudio.com/docs/cpp/config-mingw) or Linux(apt-get install build-essential and sudo apt-get install pkg-config) then install [Oracle Instant Client and SDK](https://medium.com/@utranand/how-to-connect-golang-to-oracle-on-windows-64-bit-using-go-oci8-library-ab9ed0511b20) ([Alternative Link](https://web.archive.org/web/20230105062526/https://medium.com/@utranand/how-to-connect-golang-to-oracle-on-windows-64-bit-using-go-oci8-library-ab9ed0511b20)) - Yes, even it's in Thai, it's very useful, great job Puttipong Utranand
3. Clone the repository with git
4. go to the main folder and 

        go install


## Functional Flux
You have a project with a "db" folder (it doesn't need to be named "db" at all) that needs to be run.

You just need to:
1. Go to the db folder
2. Write "beryl ah -n projectName"
3. Press Enter
4. The app will load the folders and subfolders with the files
5. A message will appear informing the number of folders and files added
6. Write "beryl u -n projectName"
7. Now the files will be executed in alphabetical order (ex.: 01-tables/001-table1.sql; 01-tables/002-table2.sql; 02-triggers/001-trigger1.sql).
8. If everything were OK, you'll receive a success message.

## Code Updates
That's why we're here, Beryl SQL Helper (BSH) is a smart boi.

You just need to:
1. Write "beryl vr -n projectName"
2. Press Enter
3. A message will appear informing the number of folders and files added and updated (BSH reads the last modification from the file).
4. Write "beryl u -n projectName"
5. BSH now have a list of added and modified files and folders, it will just run those files (in alphabetical order as above).
6. If everything were OK, you'll receive a success message.

## Database Connection
Since you've added the project to BSH, it will create a file called "c_ProjectName_.cnf" with a basic connection to localhost, use __beryl tc -n projectName__ to verify the conectivity.

## Have external variables?
No problem, just add in the main folder a file named "ProjectName.bsh" with variable";"value inside it.

Ex.:

    ${database_name};MyNewDB
    ${specific_table_name};OutstandingTable
    
## Best pratices
For files: 

    VYYYY.MM.DD.HH.MM.SS.XXXX_NAMEHERE.sql - like V2022.11.15.00.37.00.0001_CREATE_DATABASE.sql

For variables:

    ${variable_name} - like ${database_name}

We recommend putting an

    USE YOURDATABASENAME;
    
or (for PostgreSQL)

    \c YOURDATABASENAME;

as the start of the file as long as it is not a "CREATE DATABASE" file.    

## Commands

> ### ___--help / -h____
> Show this text.

> ### ___--version / -v___
> Shows the installed version of the code.

> ### ___sa___
> Shows all main folders for each project.

> ### ___s -n projectName___
> Shows the data of the selected project.

> ### ___va___
> Verifies all projects and covered folders for updates.

> ### ___vr -n projectName___
> Verifies a specific project and covered folders for updates.

> ### ___an -n projectName -l projectLocation___
> Adds a new project and its folder to the app.

> ### ___ah -n projectName___
> Adds the current folder to the app.

> ### ___ua___
> Updates all projects added to the app.

> ### ___u -n projectName___
> Updates a specific project.
> 
>> #### ___--force / -f (only with --update)___
>> (be careful) Re-run all files in all folders. 

> ### ___tc -n projectName___
> Test the connection with the server/database.

> ### ___r -i id -n projectName___
> Rename the selected project. (ID can be viewed in --showall)

> ### ___rp -n projectName -w newProjectLocation___
> Changes in the internal db map to the project folder. (THIS DOES NOT REPLACE FILES OR FOLDERS)

> ### ___d -n projectName___
> Delete in the internal db map the project. (THIS DOES NOT DELETE FILES OR FOLDERS)


## Common questions
> Q.: **Why the name Beryl?**
>> A.: https://en.wikipedia.org/wiki/Beryl#Etymology

> Q.: **Something went wrong, a error message appeared!**
>> A.: The app just mirrors the db errors, an error file was created in the main folder.

> Q.: **In what databases BSH can be used?**
>> A.: MariaDB and MySQL (basically MySQL adapter) - OracleDB and MSSQL Server still not tested.

> Q.: **Can I use it in my job/college work/company?**
>> A.: For sure, just spread the name of Beryl SQL Helper (Give us credit please).

> Q.: **How can I help the project?**
>> A.: Testing, reporting issues, giving ideas, codes, and [donating](https://ko-fi.com/mrGlasses) - optional, we don't even ask for money on company use.

> Q.: **Does BSH have an auto-update?**
>> A.: We don't think it's useful, but we can work on it.

> Q.: **Can I exclude some folders/files?**
>> A.: Not at this moment, we can work on that feature.

> Q.: **What is BSH?**
>> A.: ~~You have ADHD or at least AD, visit a psychologist~~ See __Code Updates__ above.
