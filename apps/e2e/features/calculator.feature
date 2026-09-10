Feature: Calculator web UI
  As a user
  I want to run basic arithmetic in the browser
  So that operations are executed by the Go backend

  Background:
    Given I open the calculator page

  Scenario: Empty history shows placeholder
    Then the history empty message should be "No history yet"

  Scenario Outline: Successful operations
    When I set operand A to "<a>" and B to "<b>"
    And I click the "<op>" operation
    Then the result should show "<result>"
    And no error should be visible

    Examples:
      | a  | b | op       | result |
      | 10 | 2 | add      | 12     |
      | 10 | 2 | subtract | 8      |
      | 10 | 2 | multiply | 20     |
      | 10 | 2 | divide   | 5      |

  Scenario: Division by zero shows an error
    When I set operand A to "10" and B to "0"
    And I click the "divide" operation
    Then an error containing "division by zero" should be visible
    And the result should show "—"

  Scenario: Multiply shows correct result
    When I set operand A to "4" and B to "5"
    And I click the "multiply" operation
    Then the result should show "20"
    And no error should be visible

  Scenario: Successful add appears in history
    When I set operand A to "10" and B to "2"
    And I click the "add" operation
    Then the result should show "12"
    And the history should contain "10 + 2 = 12"

  Scenario: Division by zero does not add a history row
    When I set operand A to "10" and B to "0"
    And I click the "divide" operation
    Then an error containing "division by zero" should be visible
    And the result should show "—"
    And the history should not contain "10 ÷ 0 = 0"

  Scenario: History is newest-first
    When I set operand A to "1" and B to "1"
    And I click the "add" operation
    And I set operand A to "2" and B to "2"
    And I click the "add" operation
    Then the history first entry should contain "2 + 2 = 4"
    And the history should contain "1 + 1 = 2"

  Scenario: History persists after reload
    When I set operand A to "3" and B to "4"
    And I click the "add" operation
    Then the history should contain "3 + 4 = 7"
    When I reload the calculator page
    Then the history should contain "3 + 4 = 7"

  Scenario: History is capped at 10 entries
    Given I have performed 10 successful operations
    When I set operand A to "99" and B to "1"
    And I click the "add" operation
    Then the history should have at most 10 entries
    And the history should contain "99 + 1 = 100"
