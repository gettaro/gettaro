import { useState, useEffect, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { useOrganizationStore } from '../stores/organization'
import Api from '../api/api'
import { Member } from '../types/member'
import { Title } from '../types/title'
import { Search } from 'lucide-react'

export default function Members() {
  const navigate = useNavigate()
  const { currentOrganization } = useOrganizationStore()
  
  // Data state
  const [members, setMembers] = useState<Member[]>([])
  const [titles, setTitles] = useState<Title[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')

  // Load members
  useEffect(() => {
    const loadMembers = async () => {
      if (!currentOrganization?.id) return
      try {
        setLoading(true)
        const membersData = await Api.getOrganizationMembers(currentOrganization.id)
        setMembers(membersData)
      } catch (err) {
        console.error('Error loading members:', err)
      } finally {
        setLoading(false)
      }
    }
    loadMembers()
  }, [currentOrganization?.id])

  // Load titles
  useEffect(() => {
    const loadTitles = async () => {
      if (!currentOrganization?.id) return
      try {
        const titlesData = await Api.getOrganizationTitles(currentOrganization.id)
        setTitles(titlesData)
      } catch (err) {
        console.error('Error loading titles:', err)
      }
    }
    loadTitles()
  }, [currentOrganization?.id])

  // Filter members by search query
  const filteredMembers = useMemo(() => {
    if (!searchQuery.trim()) {
      return members
    }
    
    const query = searchQuery.toLowerCase().trim()
    return members.filter(member => {
      const username = member.username?.toLowerCase() || ''
      const email = member.email?.toLowerCase() || ''
      const title = member.title_id 
        ? titles.find(t => t.id === member.title_id)?.name?.toLowerCase() || ''
        : ''
      
      return username.includes(query) || 
             email.includes(query) || 
             title.includes(query)
    })
  }, [members, titles, searchQuery])

  if (!currentOrganization) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-7xl mx-auto">
          <p className="text-muted-foreground">Please select an organization</p>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">Members</h1>
          <p className="text-muted-foreground">
            View and manage all members in your organization.
          </p>
        </div>

        {/* Search Box */}
        <div className="mb-6">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-muted-foreground" />
            <input
              type="text"
              placeholder="Search by name, email, or title..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-3 border border-border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
            />
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center items-center py-12">
            <span className="text-muted-foreground">Loading members...</span>
          </div>
        ) : filteredMembers.length === 0 ? (
          <div className="text-center py-12 bg-card rounded-lg border border-border">
            <p className="text-muted-foreground">
              {searchQuery ? 'No members found matching your search.' : 'No members found.'}
            </p>
          </div>
        ) : (
          <>
            <div className="mb-4 flex items-center justify-between">
              <p className="text-sm text-muted-foreground">
                {filteredMembers.length} {filteredMembers.length === 1 ? 'member' : 'members'}
                {searchQuery && ` found`}
              </p>
            </div>
            
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {filteredMembers.map((member) => {
                const memberTitle = member.title_id 
                  ? titles.find(t => t.id === member.title_id)?.name 
                  : null
                
                return (
                  <button
                    key={member.id}
                    onClick={() => navigate(`/members/${member.id}/profile`)}
                    className="w-full flex items-center space-x-3 p-4 bg-card rounded-lg border border-border hover:border-primary/50 hover:bg-muted/30 transition-colors text-left"
                  >
                    {/* Avatar */}
                    <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center flex-shrink-0">
                      <span className="text-primary font-medium text-base">
                        {member.username?.charAt(0).toUpperCase() || '?'}
                      </span>
                    </div>
                    
                    {/* Member Info */}
                    <div className="flex-1 min-w-0">
                      <h4 className="font-semibold text-foreground text-base mb-1 truncate">
                        {member.username || 'Unknown'}
                      </h4>
                      {memberTitle && (
                        <p className="text-sm text-muted-foreground truncate mb-1">
                          {memberTitle}
                        </p>
                      )}
                      {member.email && (
                        <p className="text-xs text-muted-foreground truncate">
                          {member.email}
                        </p>
                      )}
                    </div>
                  </button>
                )
              })}
            </div>
          </>
        )}
      </div>
    </div>
  )
}

