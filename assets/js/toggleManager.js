export function setup() {
  const toggleElements = document.querySelectorAll("[data-toggle]");

  toggleElements.forEach((element) => {
    element.addEventListener("click", () => {
      const targetSelector = element.getAttribute("data-toggle-target");
      const toggleClass = element.getAttribute("data-toggle-class");

      if (targetSelector && toggleClass) {
        const targetElement = document.querySelector(targetSelector);

        if (targetElement) {
          targetElement.classList.toggle(toggleClass);
        }
      }
    });
  });
}
