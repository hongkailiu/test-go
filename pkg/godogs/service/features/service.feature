Feature: log in
  As a user, I should be able to login using service

  Scenario: Log in with the correct username-password pair
    When I user username "username" and password "password" as input of the service
    Then the service returns "OK"

  Scenario: Log in with the wrong username-password pair
    When I user username "username" and password "wrong-password" as input of the service
    Then the service returns "NOT-OK"