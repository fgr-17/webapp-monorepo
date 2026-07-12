Feature: Calculator web UI
  As a user
  I want to run basic arithmetic in the browser
  So that operations are executed by the Go backend

  Background:
    Given I open the calculator page

  Scenario Outline: Successful operations
    When I set operand A to "<a>" and B to "<b>"
    And I click the "<op>" operation
    Then the result should show "<result>"
    And no error should be visible

    Examples:
      | a  | b | op       | result |
      | 10 | 2 | add     | 12     |
      | 10 | 2 | subtract | 8      |
      | 10 | 2 | multiply | 20     |
      | 10 | 2 | divide   | 5      |

  Scenario: Division by zero shows an error
    When I set operand A to "10" and B to "0"
    And I click the "divide" operation
    Then an error containing "division by zero" should be visible
    And the result should show "—"
