export interface ModrinthProject {
    project_id: string;
    slug: string;
    title: string;
    description: string;
    categories: string[];
    icon_url: string;
    client_side: string;
    server_side: string;
    downloads: number;
    follows: number;
    project_type: string;
    author: string;
    versions: string[];
    latest_version?: string;
    license: string;
    gallery: string[];
    date_created: string;
    date_modified: string;
}

export interface ModrinthSearchResponse {
    hits: ModrinthProject[];
    offset: number;
    limit: number;
    total_hits: number;
}
