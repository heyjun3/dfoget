'use client'
import { useFormState } from 'react-dom'
import { loginAction } from './form'
import { Button } from '@/components/ui/button'

const initialState = { message: '' }

export default function Page() {
  const [state, formAction] = useFormState(loginAction, initialState)
  return (
    <>
      <h1>Login Test</h1>
      <form action={formAction}>
        <div><input type="text" name="user_id" placeholder="User ID" /></div>
        <div><input type="password" name="password" placeholder="Password" /></div>
        <Button type='submit'>Click me</Button>
        <p>{state?.message}</p>
      </form>
    </>
  )
}
