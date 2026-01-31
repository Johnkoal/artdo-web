# 🚀 Artdotech Web Ecosystem
**Transforming Digital Ideas into Robust Solutions.**

Welcome to the official repository of the Artdotech website. This project serves as the digital gateway for a technology consultancy and software development firm based in Bogotá, Colombia.

---

## ✨ Core Features

### 🌍 Intelligent Localization
Full multi-language support (**Spanish** & **English**) with automatic language detection and manual overrides. Content managed through clean JSON locale files.

### 🎨 Premium User Experience
*   **Modern Hero Section**: High-impact visuals with dynamic gradient overlays.
*   **Zig-Zag Architecture**: Optimized layout for readability, showcasing Mission and Vision with artistic balance.
*   **The Chain of Value**: A customized services timeline that visually connects business goals with technical execution.
*   **Fully Responsive**: Fluid design that scales beautifully from 4K desktops to small mobile devices.

### 🔍 Elite SEO & Visibility
*   **Dynamic Meta Tags**: Automated Open Graph (OG), Canonical URLs, and unique meta descriptions for every page.
*   **Indexation Ready**: Integrated `sitemap.xml` and `robots.txt` served directly for search engine efficiency.

### 🛡️ Security & Reliability
*   **Anti-Spam Barrier**: Integrated Google reCAPTCHA v2 in the contact ecosystem.
*   **Backend Integrity**: Powered by **Go (Golang)** for maximum performance and memory safety.
*   **Containerization**: Ready for production with a modular **Dockerfile**.

---

## 🛠 Technology Stack

### Backend Engine
*   **Go (Golang)**: High-performance routing and template rendering.
*   **Standard Library**: Robust HTTP handling without heavy overhead.

### Frontend Aesthetics
*   **HTML5 / CSS3**: Semantic structure with custom design tokens.
*   **Bootstrap 4.5**: Grid system for reliable responsiveness.
*   **FontAwesome 6**: Extensive icon library for a modern look.

### Domain Mastery (What we do)
The site showcases our expertise in:
*   **Languages**: Go, Python, Node.js, .NET, Java.
*   **Frontend**: React, Vue.js.
*   **Mobile**: Flutter, Swift, Kotlin.
*   **Cloud**: AWS, Azure, Docker.
*   **Data**: PostgreSQL, MySQL, SQL Server.

---

## 📂 Project Structure

```text
├── locales/          # Localization JSON files (es.json, en.json)
├── static/           # CSS, Images, and SEO static files
├── templates/        # Go HTML templates (Layout, Home, Services, etc.)
├── main.go           # Server logic and route handlers
├── config.json       # App configuration (SMTP, CAPTCHA, SSL)
└── Dockerfile        # Production deployment build
```

---

## 🚀 Getting Started

### Prerequisites
*   Go 1.21 or higher.

### Local Execution
1.  Clone the repository.
2.  Configure your `config.json` with your credentials.
3.  Run the application:
    ```bash
    go run .
    ```
4.  Open `http://localhost:8080` in your browser.

### Docker Deployment
```bash
docker build -t artdo-web .
docker run -p 8080:8080 artdo-web
```

---
*Created with ❤️ by the Artdotech Team.*
