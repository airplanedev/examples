import cowsay from "cowsay";

export const hello = () => {
  return "hello!";
};

export const say = (msg) => {
  console.log(cowsay.say({ text: msg }));
};
