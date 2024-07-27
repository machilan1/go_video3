

# Maya Video with Golang
## Serving

### 1. Create .env  

Just copy the exact content in .envexample into .env

### 2. Run docker-compose  

    $ docker-compose up 

### 3. Enter psql (if psql is available on your machine)   

    $ make psql 

### 4. Create tables and seeds through sql/schema.sql and sql/seed.sql manually.


### 5. To serve the app
    $ make run

## Notices

### 1. There are three users with different roles that you can login with. Password is literally "password"

teacher@gmail.com  (recommended login with this one)   
student@gmail.com   
admin@gmail.com  


## Developing

### 1. Install Tailwind CSS

    $ make setup OS= <your os>


### 2. Tailwind CSS has to recompile every time you change the templates.

Open a new terminal and run

    $ make tailwindcss-watch 

and the tailwindcss CLI will monitor the changes and recompile stylesheets automatically.


### 3. If you want to feel even better , you can use Air for hot reloading.
    
Install air by the instructions in the url below

https://github.com/air-verse/air

There is a .air.toml in project root and you can utilize the setting.

Run

    $ air

to serve the app with hot reloading.
