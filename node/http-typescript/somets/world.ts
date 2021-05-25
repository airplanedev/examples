import cowsay from "cowsay";

export const sayButWithTypes = (msg: string) => {
  console.log(cowsay.say({ text: msg }));
};
