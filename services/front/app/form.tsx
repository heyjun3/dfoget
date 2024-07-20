'use server'
export async function loginAction(state: any, formData: FormData) {
  const user_id = formData.get('user_id')
  const password = formData.get('password')
  console.log({user_id, password})
  return {
    message: 'OK',
  }
}
