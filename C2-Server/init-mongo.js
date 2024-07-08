db = db.getSiblingDB('C2Server'); // Switch to C2Server database

db.createCollection('profile'); // Create profiles collection

db.profile.insertMany([
    {
        username: 'user1',
        password: 'password1'
    },
    {
        username: 'user2',
        password: 'password2'
    }
]);
