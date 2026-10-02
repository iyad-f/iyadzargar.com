// SPDX-FileCopyrightText: 2026 Iyad
// SPDX-License-Identifier: Apache-2.0

(function () {
	"use strict";

	function initThemeToggle() {
		const button = document.getElementById("theme-toggle");
		if (!button) {
			return;
		}

		button.addEventListener("click", function () {
			const root = document.documentElement;
			root.classList.toggle("dark");
			localStorage.theme = root.classList.contains("dark") ? "dark" : "light";
		});
	}

	function initMobileNav() {
		const toggle = document.getElementById("nav-toggle");
		const menu = document.getElementById("mobile-menu");
		if (!toggle || !menu) {
			return;
		}

		function setOpen(open) {
			menu.classList.toggle("hidden", !open);
			toggle.setAttribute("aria-expanded", String(open));
		}

		toggle.addEventListener("click", function () {
			setOpen(menu.classList.contains("hidden"));
		});

		menu.querySelectorAll("a").forEach(function (link) {
			link.addEventListener("click", function () {
				setOpen(false);
			});
		});
	}

	function initTypewriter() {
		const name = document.querySelector("[data-text]");
		if (!name) {
			return;
		}

		const full = name.dataset.text;
		const space = full.indexOf(" ");
		const split = space < 0 ? full.length : space + 1;

		function render(count) {
			name.textContent = "";
			name.appendChild(document.createTextNode(full.slice(0, Math.min(count, split))));

			if (count > split) {
				const tail = document.createElement("span");
				tail.style.color = "var(--signal)";
				tail.textContent = full.slice(split, count);
				name.appendChild(tail);
			}
		}

		if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
			render(full.length);
			return;
		}

		let count = 0;
		let direction = 1;

		function tick() {
			render(count);

			if (direction === 1) {
				count += 1;
				if (count > full.length) {
					count = full.length;
					direction = -1;
					setTimeout(tick, 2500);
					return;
				}
				setTimeout(tick, 110);
				return;
			}

			count -= 1;
			if (count < 0) {
				count = 0;
				direction = 1;
				setTimeout(tick, 700);
				return;
			}
			setTimeout(tick, 55);
		}

		tick();
	}

	function initTerminal() {
		const term = document.getElementById("term");
		if (!term) {
			return;
		}

		const dot = document.getElementById("term-dot");
		const label = document.getElementById("term-label");
		const toggle = document.getElementById("term-toggle");

		const accent = "var(--term-accent)";
		const bright = "var(--term-fg)";
		const ok = "var(--term-ok)";
		const dim = "var(--term-dim)";
		const reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

		let running = false;
		let generation = 0;
		let atBottom = true;

		function paint(text, color) {
			return `<span style="color:${color}">${text}</span>`;
		}

		function push(html) {
			const row = document.createElement("div");
			row.innerHTML = html;
			term.appendChild(row);

			while (term.children.length > 200) {
				term.removeChild(term.firstChild);
			}

			if (atBottom) {
				term.scrollTop = term.scrollHeight;
			}
		}

		function setState(on) {
			running = on;

			if (dot) {
				dot.style.background = on ? ok : dim;
			}

			if (label) {
				label.textContent = on ? "stop" : "run";
			}
		}

		const boot = [
			`level=info msg=${paint('"starting server"', bright)} addr=:8080`,
			`level=info msg=${paint('"static assets mounted"', bright)} dir=web/static`,
			`level=info msg=${paint('"routes registered"', bright)} count=4`,
			`level=info msg=${paint('"listening"', bright)} addr=:8080`,
		];

		const paths = ["/", "/projects", "/about", "/contact"];

		function request() {
			const path = paths[Math.floor(Math.random() * paths.length)];
			const duration =
				Math.random() < 0.6
					? Math.floor(80 + Math.random() * 900) + "µs"
					: (0.1 + Math.random() * 4).toFixed(1) + "ms";

			return `level=info ${paint("GET", accent)} ${path} ${paint("200", ok)} ${duration}`;
		}

		function traffic(id) {
			if (id !== generation || !running) {
				return;
			}

			push(request());
			setTimeout(function () {
				traffic(id);
			}, 1200 + Math.random() * 1800);
		}

		function bootSequence(id) {
			let line = 0;

			function step() {
				if (id !== generation) {
					return;
				}

				if (line < boot.length) {
					push(boot[line]);
					line += 1;
					setTimeout(step, reduce ? 0 : 320);
					return;
				}

				if (!reduce) {
					setTimeout(function () {
						traffic(id);
					}, 900);
				}
			}

			step();
		}

		function start() {
			if (running) {
				return;
			}

			generation += 1;
			setState(true);
			const id = generation;

			const row = document.createElement("div");
			row.innerHTML = `${paint("$", accent)} <span style="color:${bright}"></span>`;
			term.appendChild(row);

			if (atBottom) {
				term.scrollTop = term.scrollHeight;
			}

			const out = row.lastChild;
			const command = "./iyadzargar.com";

			if (reduce) {
				out.textContent = command;
				bootSequence(id);
				return;
			}

			let typed = 0;

			function type() {
				if (id !== generation) {
					return;
				}

				out.textContent = command.slice(0, typed);

				if (typed <= command.length) {
					typed += 1;
					setTimeout(type, 80);
					return;
				}

				bootSequence(id);
			}

			type();
		}

		function stop() {
			if (!running) {
				return;
			}

			generation += 1;
			setState(false);

			const lines = [
				paint("^C", dim),
				`level=warn msg=${paint('"shutdown signal received"', bright)}`,
				`level=info msg=${paint('"server is now ready to exit, bye bye..."', bright)}`,
			];

			let line = 0;

			function step() {
				if (line < lines.length) {
					push(lines[line]);
					line += 1;
					setTimeout(step, reduce ? 0 : 260);
				}
			}

			step();
		}

		term.addEventListener("scroll", function () {
			atBottom = term.scrollHeight - term.scrollTop - term.clientHeight < 24;
		});

		if (toggle) {
			toggle.addEventListener("click", function () {
				if (running) {
					stop();
				} else {
					start();
				}
			});
		}

		start();
	}

	function initProjectFilter() {
		const filters = document.getElementById("filters");
		if (!filters) {
			return;
		}

		const cards = Array.from(document.querySelectorAll("#grid > article"));

		filters.addEventListener("click", function (event) {
			const button = event.target.closest("button[data-filter]");
			if (!button) {
				return;
			}

			const filter = button.dataset.filter;

			filters.querySelectorAll("button").forEach(function (other) {
				other.dataset.active = String(other === button);
			});

			cards.forEach(function (card) {
				card.hidden = filter !== "all" && !card.dataset.tags.split(" ").includes(filter);
			});
		});
	}

	function initToasts() {
		const region = document.getElementById("toast-region");
		if (!region) {
			return;
		}

		const reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

		function dismiss(toast) {
			if (!toast) {
				return;
			}
			if (reduce) {
				toast.remove();
				return;
			}
			toast.classList.add("translate-x-4", "opacity-0");
			setTimeout(function () {
				toast.remove();
			}, 300);
		}

		// Toasts are added by htmx, so close is handled by delegation.
		region.addEventListener("click", function (event) {
			const close = event.target.closest("[data-toast-close]");
			if (close) {
				dismiss(close.closest("[data-kind]"));
			}
		});

		// Animate in and auto-dismiss each toast htmx swaps into the region.
		document.body.addEventListener("htmx:oobAfterSwap", function () {
			const toast = region.firstElementChild;
			if (!toast || toast.dataset.managed) {
				return;
			}
			toast.dataset.managed = "true";

			if (toast.dataset.kind === "ok") {
				const form = document.getElementById("contact-form");
				if (form) {
					form.reset();
				}
			}

			requestAnimationFrame(function () {
				toast.classList.remove("translate-x-4", "opacity-0");
			});

			setTimeout(function () {
				dismiss(toast);
			}, 5000);
		});
	}

	function initContactForm() {
		const form = document.getElementById("contact-form");
		if (!form) {
			return;
		}

		const button = form.querySelector("button[type=submit]");
		if (!button) {
			return;
		}

		// Swap the button to "Sending..." for the duration of the htmx request.
		form.addEventListener("htmx:beforeRequest", function () {
			button.dataset.label = button.innerHTML;
			button.textContent = "Sending...";
		});

		form.addEventListener("htmx:afterRequest", function () {
			if (button.dataset.label) {
				button.innerHTML = button.dataset.label;
				delete button.dataset.label;
			}
		});
	}

	initThemeToggle();
	initMobileNav();
	initTypewriter();
	initTerminal();
	initProjectFilter();
	initContactForm();
	initToasts();
})();
