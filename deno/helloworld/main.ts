import { parse } from "https://deno.land/std@0.85.0/flags/mod.ts"

const args = parse(Deno.args)
let { name } = args

console.log(`Hello ${name}!`)
