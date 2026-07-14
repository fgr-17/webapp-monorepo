import os
import time

from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait

from pages.calculator_page import CalculatorPage


def before_all(context):
    context.base_url = os.environ.get("BASE_URL", "http://localhost:8085").rstrip("/")
    context.selenium_url = os.environ.get(
        "SELENIUM_REMOTE_URL", "http://localhost:4444/wd/hub"
    )


def before_scenario(context, scenario):
    options = Options()
    options.add_argument("--headless=new")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--window-size=1280,800")

    # Selenium Grid can briefly reject sessions while nodes register.
    last_err = None
    for _ in range(20):
        try:
            context.driver = webdriver.Remote(
                command_executor=context.selenium_url,
                options=options,
            )
            break
        except Exception as err:  # noqa: BLE001 — retry until Grid is ready
            last_err = err
            time.sleep(1)
    else:
        raise RuntimeError(f"Could not create WebDriver session: {last_err}")

    context.driver.set_page_load_timeout(30)
    context.wait = WebDriverWait(context.driver, 10)
    context.page = CalculatorPage(context.driver, context.wait, context.base_url)


def after_scenario(context, scenario):
    driver = getattr(context, "driver", None)
    if driver is not None:
        driver.quit()
