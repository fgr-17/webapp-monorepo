from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC


class CalculatorPage:
    def __init__(self, driver, wait, base_url):
        self.driver = driver
        self.wait = wait
        self.base_url = base_url

    def open(self):
        self.driver.get(f"{self.base_url}/")
        self.wait.until(EC.presence_of_element_located((By.ID, "calc-form")))

    def _set_input(self, element_id, value):
        el = self.driver.find_element(By.ID, element_id)
        self.driver.execute_script("arguments[0].value = '';", el)
        el.send_keys(str(value))

    def set_operands(self, a, b):
        self._set_input("a", a)
        self._set_input("b", b)

    def click_operation(self, op):
        button = self.wait.until(
            EC.element_to_be_clickable((By.CSS_SELECTOR, f'button[data-op="{op}"]'))
        )
        button.click()

    def result_text(self):
        # Wait until the UI leaves the idle state after a click.
        self.wait.until(
            lambda d: d.find_element(By.ID, "result").text != "Result: —"
            or self.error_visible()
        )
        return self.driver.find_element(By.ID, "result").text

    def error_visible(self):
        el = self.driver.find_element(By.ID, "error")
        return el.is_displayed() and bool(el.text.strip())

    def error_text(self):
        return self.driver.find_element(By.ID, "error").text
