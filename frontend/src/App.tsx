import React, { useState } from 'react'
import { Transcribe, FFmpegVersion } from '../wailsjs/go/backend/App'

export default function App() {
    const [path, setPath] = useState('')
    const [log, setLog] = useState('')

    const transcribe = async () => {
        setLog('Transcribiendo...')
        try {
            const srtPath = await Transcribe(path)
            setLog(`Listo → ${srtPath}`)
        } catch (e: any) {
            setLog('Error: ' + e?.message)
        }
    }

    const ffver = async () => {
        try {
            const v = await FFmpegVersion()
            setLog(v)
        } catch (e: any) {
            setLog('FFmpeg no encontrado')
        }
    }

    return (
        <div style={{ fontFamily: 'Inter, system-ui, sans-serif', padding: 16 }}>
            <h1>Montaje</h1>
            <p>Skeleton con transcripción a SRT usando faster‑whisper.</p>
            <div style={{ display: 'flex', gap: 8 }}>
                <input style={{ flex: 1 }} placeholder="Ruta a tu video/audio" value={path} onChange={e => setPath(e.target.value)} />
                <button onClick={transcribe}>Transcribir</button>
                <button onClick={ffver}>FFmpeg?</button>
            </div>
            <pre style={{ marginTop: 16, background: '#111', color: '#0f0', padding: 12, borderRadius: 8 }}>{log}</pre>
        </div>
    )
}