# InvoiceLLM Backend
## Swipe Hiring Challenge

This is the backend for the website where you can upload an invoice and get customer details, product details and total amount in json format. You can also choose to store the invoices and retrieve it, must be logged in to access these features.

You can use the pdfs provided in the sampledata folder to test the project

## DataBase Schema

### users (postgres)
    id        : UUID
    user_name : TEXT
    passwd    : TEXT

### invoices (mongoDB)
    _id   : ObjectId
    userid: string
    invoice: Object
        customerdetails: Object
            name: string
            address
            phonenumber
            email
        productdetails: Array
            Object
                name
                rate
                quantity
                totalamount
    totalamount
    amountpayable

## API Endpoints

### General APIs: (without auth)
    GET /          : returns hello user, no body and headers required
    POST /upload   : takes pdf file and returns the invoice extracted from pdf in json format,

### Auth APIs:
    POST /signup   : crates a new user, takes username and password in request body
    POST /login    : returns a new jwt token, takes username and password in request body
    
### Restricted: (with auth)
    POST /invoices :  adds the extracted invoice into the database, requires the extracted invoice data as request body
    GET /invoices  : returns all the invoices saved by a particular user

Note: Jwt token expiry after logout and refresh token not implemented.