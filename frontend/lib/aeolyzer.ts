export type AccountType = "brand" | "agency"
export type Reach =
  | "global"
  | "primary-market"
  | "nationwide"
  | "regional"
  | "local"

export interface ProjectProfile {
  session_id: string
  account_type: AccountType
  domain: string
  brand_name: string
  reach: Reach
  country_code: string
  country_name: string
  language: string
  competitors: string[]
  icon_url?: string
  description?: string
  category?: string
  prompts: string[]
}

export interface SiteInspection {
  canonical_url: string
  suggested_brand_name: string
  description: string
  category: string
  icon_url: string
  competitor_candidates: string[]
}

export interface DashboardFrame {
  schema_version: string
  surface: "audit_dashboard"
  generated_at: string
  project: {
    brand_name: string
    domain: string
    location: string
    language: string
  }
  tabs: Array<{ id: string; label: string; enabled: boolean }>
  prompts: string[]
}

export interface ChatMessage {
  id: string
  role: "user" | "assistant"
  content: string
}

export interface CountryOption {
  code: string
  name: string
  languageCode: string
}

export interface LanguageOption {
  code: string
  name: string
}

// Strip trailing slash statically to guarantee predictable URL construction downstream and prevent SSR route mismatches.
export const API_URL =
  process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") ?? "http://localhost:8080"

const countryCodes = `
AD AE AF AG AI AL AM AO AQ AR AS AT AU AW AX AZ BA BB BD BE BF BG BH BI BJ BL BM BN
BO BQ BR BS BT BV BW BY BZ CA CC CD CF CG CH CI CK CL CM CN CO CR CU CV CW CX CY CZ
DE DJ DK DM DO DZ EC EE EG EH ER ES ET FI FJ FK FM FO FR GA GB GD GE GF GG GH GI GL
GM GN GP GQ GR GS GT GU GW GY HK HM HN HR HT HU ID IE IL IM IN IO IQ IR IS IT JE JM
JO JP KE KG KH KI KM KN KP KR KW KY KZ LA LB LC LI LK LR LS LT LU LV LY MA MC MD ME
MF MG MH MK ML MM MN MO MP MQ MR MS MT MU MV MW MX MY MZ NA NC NE NF NG NI NL NO NP
NR NU NZ OM PA PE PF PG PH PK PL PM PN PR PS PT PW PY QA RE RO RS RU RW SA SB SC SD
SE SG SH SI SJ SK SL SM SN SO SR SS ST SV SX SY SZ TC TD TF TG TH TJ TK TL TM TN TO
TR TT TV TW TZ UA UG UM US UY UZ VA VC VE VG VI VN VU WF WS YE YT ZA ZM ZW
`.trim().split(/\s+/)

const languageCodes = `
aa ab ae af ak am an ar as av ay az ba be bg bh bi bm bn bo br bs ca ce ch co cr cs cu
cv cy da de dv dz ee el en eo es et eu fa ff fi fj fo fr fy ga gd gl gn gu gv ha he hi
ho hr ht hu hy hz ia id ie ig ii ik io is it iu ja jv ka kg ki kj kk kl km kn ko kr ks
ku kv kw ky la lb lg li ln lo lt lu lv mg mh mi mk ml mn mr ms mt my na nb nd ne ng nl
nn no nr nv ny oc oj om or os pa pi pl ps pt qu rm rn ro ru rw sa sc sd se sg si sk sl
sm sn so sq sr ss st su sv sw ta te tg th ti tk tl tn to tr ts tt tw ty ug uk ur uz ve
vi vo wa wo xh yi yo za zh zu
`.trim().split(/\s+/)

const defaultLanguageByCountry: Record<string, string> = {
  AE: "ar",
  AR: "es",
  AT: "de",
  AU: "en-AU",
  BD: "bn",
  BE: "nl",
  BR: "pt-BR",
  CA: "en-CA",
  CH: "de",
  CL: "es",
  CN: "zh-Hans",
  CO: "es",
  CZ: "cs",
  DE: "de",
  DK: "da",
  DZ: "ar",
  EG: "ar",
  ES: "es",
  FI: "fi",
  FR: "fr",
  GB: "en-GB",
  GR: "el",
  HK: "zh-Hant",
  HU: "hu",
  ID: "id",
  IE: "en-IE",
  IL: "he",
  IN: "hi",
  IR: "fa",
  IT: "it",
  JP: "ja",
  KE: "sw",
  KR: "ko",
  LK: "si",
  MA: "ar",
  MX: "es",
  MY: "ms",
  NG: "en",
  NL: "nl",
  NO: "nb",
  NZ: "en-NZ",
  PH: "tl",
  PK: "ur",
  PL: "pl",
  PT: "pt",
  RO: "ro",
  RU: "ru",
  SA: "ar",
  SE: "sv",
  SG: "en",
  TH: "th",
  TR: "tr",
  TW: "zh-Hant",
  UA: "uk",
  US: "en-US",
  VN: "vi",
  ZA: "en",
}

const languageVariants = [
  "en-AU",
  "en-CA",
  "en-GB",
  "en-IE",
  "en-NZ",
  "en-US",
  "pt-BR",
  "zh-Hans",
  "zh-Hant",
]

export function countries(): CountryOption[] {
  // Intl.DisplayNames defers heavy localization data loading to the browser runtime rather than bloating the bundle.
  const names = new Intl.DisplayNames(["en"], { type: "region" })
  return countryCodes
    .map((code) => ({
      code,
      name: names.of(code) ?? code,
      languageCode: defaultLanguageByCountry[code] ?? "en",
    }))
    .sort((a, b) => a.name.localeCompare(b.name))
}

export function languages(): LanguageOption[] {
  const names = new Intl.DisplayNames(["en"], { type: "language" })
  return [...languageCodes, ...languageVariants]
    .filter((code, index, values) => values.indexOf(code) === index)
    .map((code) => ({
      code,
      name: languageLabel(code, names),
    }))
    .sort((a, b) => a.name.localeCompare(b.name))
}

function languageLabel(code: string, names: Intl.DisplayNames): string {
  const labels: Record<string, string> = {
    "en-AU": "English (Australia)",
    "en-CA": "English (Canada)",
    "en-GB": "English (UK)",
    "en-IE": "English (Ireland)",
    "en-NZ": "English (New Zealand)",
    "en-US": "English (US)",
    "pt-BR": "Portuguese (Brazil)",
    "zh-Hans": "Chinese (Simplified)",
    "zh-Hant": "Chinese (Traditional)",
  }
  return labels[code] ?? names.of(code) ?? code
}

export function fallbackBrandName(rawURL: string): string {
  // Aggressive exception handling here prevents malformed user-input URLs from crashing the entire rendering cycle.
  try {
    const value = rawURL.includes("://") ? rawURL : `https://${rawURL}`
    const host = new URL(value).hostname.replace(/^www\./, "")
    const first = host.split(".")[0]
    return first
      .split(/[-_]/)
      .filter(Boolean)
      .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
      .join(" ")
  } catch {
    return ""
  }
}

export async function inspectSite(
  sessionID: string,
  url: string,
): Promise<SiteInspection> {
  const response = await fetch(`${API_URL}/v1/onboarding/inspect`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: sessionID, url }),
  })
  return readResponse<SiteInspection>(response)
}

export async function completeOnboarding(
  profile: Omit<ProjectProfile, "prompts">,
): Promise<DashboardFrame> {
  const category = profile.category?.trim()
  const endpoint = new URL(`${API_URL}/v1/onboarding/complete`)
  if (category) endpoint.searchParams.set("category", category)

  const input = {
    session_id: profile.session_id,
    account_type: profile.account_type,
    domain: profile.domain,
    brand_name: profile.brand_name,
    reach: profile.reach,
    country_code: profile.country_code,
    country_name: profile.country_name,
    language: profile.language,
    competitors: profile.competitors,
  }
  const response = await fetch(endpoint, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  })
  return readResponse<DashboardFrame>(response)
}

async function readResponse<T>(response: Response): Promise<T> {
  // Standardizes error throwing by unwrapping JSON errors before they hit component boundaries.
  if (response.ok) return (await response.json()) as T

  const payload = (await response.json().catch(() => null)) as {
    error?: { message?: string }
  } | null
  throw new Error(payload?.error?.message ?? "The request could not be completed.")
}

export const reachOptions: Array<{
  id: Reach
  title: string
  description: string
}> = [
  {
    id: "global",
    title: "Global",
    description: "Your brand serves audiences worldwide.",
  },
  {
    id: "primary-market",
    title: "Primary market + international",
    description: "One main market with reach beyond it.",
  },
  {
    id: "nationwide",
    title: "Nationwide",
    description: "Your brand operates across one country.",
  },
  {
    id: "regional",
    title: "Regional",
    description: "Your brand focuses on a state or region.",
  },
  {
    id: "local",
    title: "Local",
    description: "Your brand serves a city or nearby area.",
  },
]

export const quickActions = [
  {
    title: "Optimize Meta tags",
    prompt:
      "Review the meta titles and descriptions for my key pages. Suggest optimized versions that improve click-through rates (CTR) and include relevant keywords without keyword stuffing.",
  },
  {
    title: "Generate a llm.txt",
    prompt:
      "Generate an llms.txt file for my website to help AI agents understand my content better. Include key pages, sitemap location, and a concise description of my site's purpose.",
  },
  {
    title: "Review Technical SEO",
    prompt:
      "Perform a technical SEO review. Check for crawl errors, broken links, duplicate content, canonical tag issues, and mobile-friendliness problems.",
  },
  {
    title: "Make me a Task List",
    prompt:
      "Generate a prioritized SEO task list for my website. Focus on high-impact, low-effort actions I can take this week to improve my search rankings. Categorize them by Technical, Content, and Authority.",
  },
  {
    title: "Summarize recent Performance",
    prompt:
      "Summarize my website's search performance over the last 30 days. Highlight significant changes in traffic, rankings, and impressions. What went well and what needs attention?",
  },
  {
    title: "Find Keyword opportunities",
    prompt:
      "Find underutilized keyword opportunities for my niche. Look for long-tail keywords with decent search volume and low competition that I can target with new content.",
  },
  {
    title: "Quick Visibility Check",
    prompt:
      "Perform a comprehensive visibility check for my domain. Analyze my ranking for top keywords, identify visibility trends, and summarize my overall presence in search results.",
  },
  {
    title: "Audit my Homepage",
    prompt:
      "Conduct a detailed audit of my homepage. Check for on-page SEO issues, technical errors, user experience friction, and conversion optimization opportunities. Provide actionable recommendations.",
  },
  {
    title: "Analyze competitors",
    prompt:
      "Analyze my top 3 competitors. Compare their search visibility, top-performing content, and backlink profiles to mine. Highlight their strengths and my opportunities to outperform them.",
  },
  {
    title: "Draft a Blog Post",
    prompt:
      "Draft a high-quality, SEO-optimized blog post about a trending topic in my industry. Include a catchy title, headers, and a structure that targets user intent.",
  },
  {
    title: "Find content gaps",
    prompt:
      "Identify content gaps on my website compared to my top competitors. What topics are they covering that I am missing? Suggest 5 new article ideas to fill these gaps.",
  },
]

export const contentTypes = [
  {
    id: "article",
    label: "Article",
    description: "Authoritative, evidence-backed long-form analysis.",
    details: "800–3,000 words · full research · guides, comparisons and analysis",
  },
  {
    id: "blog-post",
    label: "Blog Post",
    description: "Conversational posts with one clear, useful thesis.",
    details: "600–1,500 words · standard research · opinion, how-to and case studies",
  },
  {
    id: "linkedin",
    label: "LinkedIn",
    description: "Professional posts that sound direct and human.",
    details: "150–400 words · light research · insights, stories and announcements",
  },
  {
    id: "youtube-desc",
    label: "YouTube Desc",
    description: "Benefit-led video descriptions that earn the click.",
    details: "150–300 words · hooks, timestamps, links and hashtags",
  },
  {
    id: "product-release",
    label: "Product Release",
    description: "Grounded, benefit-first launch communication.",
    details: "400–800 words · feature launches, updates and partnerships",
  },
]
