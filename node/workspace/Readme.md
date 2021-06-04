
	Example showing how you can use Airplane with [workspaces](https://docs.npmjs.com/cli/v7/using-npm/workspaces).

#### Usage

	1st create a new node.js task and initialize it.

	```bash
	$ airplane init --slug mytask ./server/tasks/hey/index.js
	```

	The task will be created and now you can deploy it.

	```bash
	$ airplane deploy ./server/tasks/hey/index.js
	```

	Deploy will find the nearest package.json `./server/package.json`, it will notice that
	the package.json defines a `airplane.root=".."`, so instead of treaing `./server` as the root
	it will treat `.` as the root directory.

	Execute the task:

	```bash
	airplane execute ./server/tasks/hey/index.js
	```
