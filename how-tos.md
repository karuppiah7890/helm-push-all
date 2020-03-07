# How Tos

These are some How Tos for the project. I want to follow them for this project
1. Write tests as much as possible. This is because I want the software to be robust
2. Reuse existing code / libraries / software whenever possible to speed up development,
and to stand on the shoulders of other giants as a lot of time and effort has been put to
make them, which doesn't have to be necessarily reinvented
3. Since we are already planning to write tests, do test driven development (TDD) whenever
possible
4. For any bug reported, write a failing test first for the bug, then fix the bug to pass the test
5. Keep the software as simple as possible - not much complexity in terms of code written,
design and architecture. To do this, write as little code as possible. Unnecessary code leads to
complexity and maintenance. No code means no maintenance and no complexity :P
6. Users are the core of the software. The software should be easy to use. If user can't use it
or understand how to use it, see behvaiour and try to change the software accordingly!
7. Everything must be automated as much as possible. 
8. And human errors can occur - that's okay. Put processes and automations in place to make sure
errors don't occur, or we catch errors fast! And if we are catching errors - make sure errors don't impact
us much or best to catch it before it impacts us - like failing release CI pipeline instead of releasing
something wrong