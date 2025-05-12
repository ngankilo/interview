
### 1. Interview Question
**"Since you’ve worked on a hybrid application with Flutter starting from September 2021, can you describe the most complex feature you implemented in that project? How did you approach it?"**  
(Kể từ khi bạn bắt đầu làm việc với ứng dụng hybrid bằng Flutter từ tháng 9/2021, bạn có thể mô tả tính năng phức tạp nhất mà bạn đã triển khai trong dự án đó không? Bạn đã tiếp cận nó như thế nào?)

#### Sample Answer
"In that project, I implemented a real-time chat feature. It was complex because it required integrating WebSocket for live messaging, handling offline states, and ensuring smooth UI updates. I used the BLoC pattern to manage the chat state, created a service class to handle WebSocket connections, and stored offline messages locally with `sqflite`. I also added retry logic for reconnection and used Flutter’s `ListView.builder` to optimize the chat list rendering."

#### Follow-up Questions & Guidance
- **Follow-up 1: "How did you ensure the WebSocket connection was reliable across different network conditions?"**  
  *Guidance*: Expect details on error handling, reconnection strategies, or libraries used.  
  *Sample Response*: "I implemented a reconnection mechanism with exponential backoff using the `web_socket_channel` package. I also added a heartbeat signal to detect connection drops and showed a loading indicator in the UI during reconnect attempts."

- **Follow-up 2: "What challenges did you face with offline support, and how did you solve them?"**  
  *Guidance*: Look for practical solutions like local storage or syncing logic.  
  *Sample Response*: "The main challenge was syncing messages once the app went online. I stored messages in `sqflite` with a 'pending' status, then used a background sync process to send them when the connection was restored, updating the UI via BLoC events."

---

### 2. Interview Question
**"You worked in a team of 9 people. How did you collaborate with your team to ensure the Flutter codebase was maintainable and scalable?"**  
(Bạn đã làm việc trong một nhóm 9 người. Bạn đã cộng tác với nhóm như thế nào để đảm bảo mã nguồn Flutter dễ bảo trì và mở rộng?)

#### Sample Answer
"Our team used a modular architecture with BLoC for state management. I collaborated by setting up clear folder structures like `screens`, `blocs`, and `services`, and we agreed on naming conventions. We used Git for version control with feature branches and conducted code reviews via pull requests. I also wrote documentation for key components and added unit tests with `flutter_test` to ensure code quality."

#### Follow-up Questions & Guidance
- **Follow-up 1: "What was your process for handling disagreements during code reviews?"**  
  *Guidance*: Look for conflict resolution skills and a focus on code quality.  
  *Sample Response*: "If we disagreed, I’d explain my reasoning with examples or documentation, like why a certain pattern was more testable. We’d discuss trade-offs and sometimes pair-program to find a solution, ensuring we aligned on the project’s goals."

- **Follow-up 2: "How did you ensure consistency in the UI across a team of 9 developers?"**  
  *Guidance*: Expect mention of design systems, reusable widgets, or tools like Figma.  
  *Sample Response*: "We created a shared `widgets` library with reusable components like buttons and text fields, styled according to a Figma design system provided by our UI/UX team. Regular sync-ups ensured everyone followed the same standards."

---

### 3. Interview Question
**"Flutter is a hybrid framework. How do you handle platform-specific features (e.g., iOS vs. Android) in your projects?"**  
(Flutter là một framework hybrid. Bạn xử lý các tính năng đặc thù của từng nền tảng (ví dụ: iOS và Android) trong dự án của mình như thế nào?)

#### Sample Answer
"For platform-specific features, I use the `Platform` class to detect the OS and write conditional logic. For example, in one project, I implemented native push notifications using `firebase_messaging`. On Android, I configured custom notification channels, while on iOS, I handled APNs setup via a native bridge with Method Channels. I kept the core logic in Flutter and isolated platform code in separate files."

#### Follow-up Questions & Guidance
- **Follow-up 1: "Can you give an example of a time you used Method Channels for a specific feature?"**  
  *Guidance*: Look for practical experience with native integration.  
  *Sample Response*: "I used Method Channels to access the device’s battery level. I wrote Swift code for iOS and Kotlin for Android, then called it from Flutter to display the battery percentage in the app."

- **Follow-up 2: "How do you test platform-specific code to ensure it works on both iOS and Android?"**  
  *Guidance*: Expect a testing strategy involving emulators, real devices, or CI/CD.  
  *Sample Response*: "I test on both iOS and Android emulators first, then use real devices for edge cases. I also set up a CI pipeline with GitHub Actions to run integration tests on both platforms whenever we push code."

---

### 4. Interview Question
**"Given your experience since 2021, how do you stay updated with Flutter’s latest features and best practices?"**  
(Với kinh nghiệm từ năm 2021, bạn làm thế nào để cập nhật với các tính năng mới nhất và thực hành tốt nhất của Flutter?)

#### Sample Answer
"I follow the official Flutter blog and changelog on GitHub to track updates like new widgets or performance improvements. I also read articles on Medium, watch Flutter talks on YouTube, and experiment with new features in side projects. For example, after null safety was introduced, I refactored a project to adopt it fully, which improved code reliability."

#### Follow-up Questions & Guidance
- **Follow-up 1: "What’s the most recent Flutter feature you’ve used, and how did it benefit your project?"**  
  *Guidance*: Look for awareness of updates like Impeller or Material 3.  
  *Sample Response*: "I recently used the Impeller rendering engine in a project after enabling it in Flutter 3.10. It reduced animation jank on iOS, giving us smoother transitions without changing the code."

- **Follow-up 2: "How do you convince your team to adopt a new Flutter feature or practice?"**  
  *Guidance*: Expect communication and demonstration skills.  
  *Sample Response*: "I’d present a small proof-of-concept showing the benefits, like faster builds or better performance, and compare it to our current approach. I’d also share documentation and discuss it in a team meeting to get buy-in."
