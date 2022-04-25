// wsRequest-Dashboard.js created with Cypress
//
// Start writing your Cypress tests below!
// If you're unfamiliar with how Cypress works,
// check out the link below and learn how to write your first test:
// https://on.cypress.io/writing-first-test

describe('Load the Endpoints', () => {
  it('Should load the Dashboard page', () => {
    cy.visit("http://localhost:3000/")
  })

  it('Should load the Process List page', () => {
    cy.visit("http://localhost:3000/process_list")
  })

  it('Should load the CPU page', () => {
     cy.visit("http://localhost:3000/cpu")
  })
});
