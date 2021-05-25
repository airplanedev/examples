// airplane.dev/slug/hello_world [generated: do not edit]

import { hello, say } from "./somejs/hello";
import { sayButWithTypes } from "./somets/world";

type Parameters = {
  name: string;
};

export default function (parameters: Parameters) {
  console.log(`Hello, ${parameters.name}!`);
  hello();
  say("hello from js land!");

  sayButWithTypes("hello from TS!");
}
