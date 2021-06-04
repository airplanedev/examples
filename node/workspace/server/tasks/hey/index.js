// Linked to Airplane task, do not edit:
// https://app.airplane.so:5000/t/test_js

const shortid = require('@team/shared');

export default async function(params){
  console.log('parameters: ', params);
  console.log('shortid', shortid());
}
