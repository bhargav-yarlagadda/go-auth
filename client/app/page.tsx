'use client'

import { userContext } from '@/components/UserWrapper'
import React, { useContext } from 'react'

const page = () => {
  const ctx = useContext(userContext)
  if(!ctx){
    return 
  }
  const {logout,user}=ctx
  return (
    <div className='flex flex-col'> 
      <span>hello {user?.email}</span>
      <button onClick={logout} className='bg-red-500 rounded-md w-44 mx-auto' >  logout</button>
    </div>
  )
}

export default page