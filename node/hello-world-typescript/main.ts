// airplane.dev/slug/hello_world [generated: do not edit]

import airplane from "airplane";

type Parameters = {
  name: string;
};

export default async function (parameters: Parameters) {
  airplane.output(`Hello, ${parameters.name}!`);
}
