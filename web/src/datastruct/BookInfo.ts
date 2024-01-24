interface ReadingPos {
    Link: string
    Percentage: string
}

interface BookInfo {
    Id: string
    Name: string
    Path: string
    CoverPath: string
    Opf: string
    ReadingPos: ReadingPos
}

export default BookInfo