# Real Time 

- Фронтенд
- Бэкенд 
- Чат 

## Forum Service

    1. Environment
        - configurations
        - linter
        - logger
        - autobuilding
        - CI (continue integration)
    2. Base logic
        - Main logic
        - App logic
    3. Layers
        - Transport (http server, router, handlers)
        - service (bisnes logic, main logic of the service - forum)
        - domen, core ( objeects, stuctures)
        - repository (main repository - conn and disconnection to the sqlite, CRUD AND OTHER DATA MANIPULATIONS)
    4. Features
        a. auth, User (Admin, Moderator, client)
            - identification (registration -> ) Sign up
                - validation, filtration
            - authontification - Sign in
            - authorization (administration)
            - logout - Sign out
        b. Post ( date and time )
            - Create 
            - Update
            - Read
            - Delete
        c. Comment, (CRUD - basic operations, specefic operations )
            - must be parent post
            - create
            - update
            - read
            - delete
        e. Additional packages (logics)
            - sentinel errors
            - filtration
            - validation
            - encryting and decrypting
            - UUID


## Chat Service

## SPA (Front-end)


