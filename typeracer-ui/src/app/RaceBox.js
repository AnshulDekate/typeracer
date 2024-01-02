
export const RaceBox = ({pos}) => {
    return (
    <div>
        <div className='race-box'>
            <div className='actual-run'>
            <div className="car" style={{ position: 'relative', left: `${(pos)}%` }}>o^^^o</div>
            </div>
            <div className='invisible'>o^^^o</div>
        </div>
    </div>
    )
}
