


Serving

1. Create .env
    Just copy the exact content in .envexample into .env


2. Run docker-compose

    $ docker-compose up 

3. Enter psql (if psql is available on your machine) 

    $ make psql 

4. Create tables and seeds through sql/schema.sql and sql/seed.sql manually.


5. To serve the app, run

    $ make run



Developing

1. Install Tailwind CSS

    $ make setup OS= <your os>


2. Tailwind CSS has to recompile every time you change the templates.

    $ make tailwindcss-watch 

    and the tailwindcss CLI will monitor the changes and recompile automatically.


3. If you want to feel even better , you can use Air if you want to utilize hot reloading.
    
    Instal air by the instructions in the url below
    
    https://github.com/air-verse/air

    There is a .air.toml in project root and you can utilize the setting.


    