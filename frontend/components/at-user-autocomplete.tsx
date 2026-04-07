"use client"

import React, { useState, useEffect, useCallback } from 'react'
import { commentAPI } from '@/lib/api'

interface AtUserAutocompleteProps {
  textareaRef: React.RefObject<HTMLTextAreaElement>
  value: string
  onChange: (val: string) => void
}

interface User {
  username: string
  avatar: string
}

export function AtUserAutocomplete({ textareaRef, value, onChange }: AtUserAutocompleteProps) {
  const [showDropdown, setShowDropdown] = useState(false)
  const [users, setUsers] = useState<User[]>([])
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [atPosition, setAtPosition] = useState<number | null>(null)
  const [searchTerm, setSearchTerm] = useState('')

  // 检测 @ 模式
  const detectAtPattern = useCallback(() => {
    const textarea = textareaRef.current
    if (!textarea) return

    const cursorPos = textarea.selectionStart
    const textBeforeCursor = value.substring(0, cursorPos)

    // 查找光标前最近的 @
    const lastAtIndex = textBeforeCursor.lastIndexOf('@')
    if (lastAtIndex === -1) {
      setShowDropdown(false)
      setAtPosition(null)
      return
    }

    // 检查 @ 后是否有空格（如果有则不触发）
    const textAfterAt = textBeforeCursor.substring(lastAtIndex + 1)
    if (textAfterAt.includes(' ') || textAfterAt.includes('\n')) {
      setShowDropdown(false)
      setAtPosition(null)
      return
    }

    const term = textAfterAt.trim()
    setAtPosition(lastAtIndex)
    setSearchTerm(term)
    setSelectedIndex(0)
  }, [value, textareaRef])

  // 监听输入变化
  useEffect(() => {
    detectAtPattern()
  }, [value, detectAtPattern])

  // Debounce 搜索
  useEffect(() => {
    if (atPosition === null) return

    const timer = setTimeout(async () => {
      if (searchTerm.length >= 1) {
        const results = await commentAPI.getAtUsers(searchTerm)
        setUsers(results)
        setShowDropdown(results.length > 0)
      }
    }, 300)

    return () => clearTimeout(timer)
  }, [searchTerm, atPosition])

  // 插入用户名
  const insertUsername = (username: string) => {
    const textarea = textareaRef.current
    if (!textarea || atPosition === null) return

    const cursorPos = textarea.selectionStart
    const before = value.substring(0, atPosition)
    const after = value.substring(cursorPos)

    const newValue = `${before}@${username} ${after}`
    onChange(newValue)
    setShowDropdown(false)
    setAtPosition(null)

    // 设置光标位置
    setTimeout(() => {
      const newCursorPos = atPosition + username.length + 2
      textarea.setSelectionRange(newCursorPos, newCursorPos)
      textarea.focus()
    }, 0)
  }

  // 键盘导航
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showDropdown) return

    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setSelectedIndex(prev => Math.min(prev + 1, users.length - 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setSelectedIndex(prev => Math.max(prev - 1, 0))
    } else if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      if (users[selectedIndex]) {
        insertUsername(users[selectedIndex].username)
      }
    } else if (e.key === 'Escape') {
      setShowDropdown(false)
    }
  }

  return (
    <div className="relative">
      <textarea
        ref={textareaRef}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={handleKeyDown}
        className="w-full min-h-[100px] rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
        rows={4}
      />
      {showDropdown && users.length > 0 && (
        <div className="absolute z-50 mt-1 w-full max-w-xs rounded-md border border-border bg-popover shadow-lg">
          {users.map((user, index) => (
            <button
              key={user.username}
              onClick={() => insertUsername(user.username)}
              className={`flex w-full items-center gap-2 px-3 py-2 text-left text-sm transition-colors ${
                index === selectedIndex
                  ? 'bg-accent text-accent-foreground'
                  : 'hover:bg-accent/50'
              }`}
            >
              <img
                src={user.avatar}
                alt={user.username}
                className="size-6 rounded-full"
              />
              <span>{user.username}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
