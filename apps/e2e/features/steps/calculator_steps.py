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


@then('the history should contain "{snippet}"')
def assert_history_contains(context, snippet):
    context.page.wait_for_history_containing(snippet)
    texts = context.page.history_texts()
    assert any(snippet in text for text in texts), (
        f"history={texts!r}, missing {snippet!r}"
    )


@then('the history should not contain "{snippet}"')
def assert_history_not_contains(context, snippet):
    # Give the UI a moment; failed ops must not paint a history row.
    context.page.wait_briefly()
    texts = context.page.history_texts()
    assert all(snippet not in text for text in texts), (
        f"history={texts!r}, unexpectedly contains {snippet!r}"
    )


@then("the history section should be visible")
def assert_history_section_visible(context):
    assert context.page.history_section_visible(), "expected history section"

@then('the history empty message should be "{text}"')
def assert_history_empty_message(context, text):
    context.page.wait_for_history_empty()
    actual = context.page.history_empty_text()
    assert actual == text, f"empty message={actual!r}, want {text!r}"


@then('the history first entry should contain "{snippet}"')
def assert_history_first_contains(context, snippet):
    context.page.wait_for_history_containing(snippet)
    first = context.page.history_first_text()
    assert snippet in first, f"first entry={first!r}, missing {snippet!r}"


@when("I reload the calculator page")
def reload_calculator(context):
    context.page.reload()


@given("I have performed 10 successful operations")
def perform_ten_operations(context):
    for i in range(10):
        context.page.set_operands(str(i), "1")
        context.page.click_operation("add")
        context.page.wait_for_history_containing(f"{i} + 1 = {i + 1}")


@then("the history should have at most 10 entries")
def assert_history_cap(context):
    count = context.page.history_count()
    assert count <= 10, f"history has {count} entries, want <= 10"    
