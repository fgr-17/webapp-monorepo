from behave import given, then, when


@given("I open the calculator page")
def open_calculator(context):
    context.page.open()


@when('I set operand A to "{a}" and B to "{b}"')
def set_operands(context, a, b):
    context.page.set_operands(a, b)


@when('I click the "{op}" operation')
def click_operation(context, op):
    context.page.click_operation(op)


@then('the result should show "{result}"')
def assert_result(context, result):
    text = context.page.result_text()
    expected = f"Result: {result}"
    assert text == expected, f"result text={text!r}, want {expected!r}"


@then("no error should be visible")
def assert_no_error(context):
    assert not context.page.error_visible(), (
        f"unexpected error: {context.page.error_text()!r}"
    )


@then('an error containing "{snippet}" should be visible')
def assert_error_contains(context, snippet):
    assert context.page.error_visible(), "expected error to be visible"
    text = context.page.error_text()
    assert snippet in text, f"error text={text!r}, missing {snippet!r}"
